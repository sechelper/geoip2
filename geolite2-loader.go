package geoip

import (
	"crypto/sha256"
	"database/sql"
	"encoding/csv"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/sechelper/geoip2/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var tmpDir = os.TempDir()

const (
	downloadUrl = "https://download.maxmind.com/app/geoip_download?edition_id=%s&license_key=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&suffix=%s"

	asnBlocksIPv4FilePrefix     = "GeoLite2-ASN-Blocks-IPv4"
	cityBlocksIPv4FilePrefix    = "GeoLite2-City-Blocks-IPv4"
	cityLocationsFilePrefix     = "GeoLite2-City-Locations"
	countryIPv4BlocksFilePrefix = "GeoLite2-Country-Blocks-IPv4"
	countryLocationsFilePrefix  = "GeoLite2-Country-Locations"
)

var Languages = []string{"de", "en", "es", "fr", "ja", "pt-BR", "ru", "zh-CN"}

type GeoLite2Loader struct {
	db *sql.DB
}

func NewGeoLite2Loader(db *sql.DB) *GeoLite2Loader {
	return &GeoLite2Loader{
		db: db,
	}
}

func (loader *GeoLite2Loader) loading(asn, city, country string) error {

	log.Debug().Msg("开始加载 [ASNBlocksIPv4] " + filepath.Join(asn, asnBlocksIPv4FilePrefix+".csv"))
	if err := loader.loadASNBlocksIPv4Csv(filepath.Join(asn, asnBlocksIPv4FilePrefix+".csv"),
		asnBlocksIPv4Sql); err != nil {
		return err
	}

	log.Debug().Msg("开始加载 [CityBlocksIPv4] " + filepath.Join(city, cityBlocksIPv4FilePrefix+".csv"))
	if err := loader.loadCityBlocksIPv4Csv(filepath.Join(city, cityBlocksIPv4FilePrefix+".csv"),
		cityBlocksIPv4Sql); err != nil {
		return err
	}

	for _, language := range Languages {
		log.Debug().Msg("开始加载 [CityLocations-" + language + "] " + filepath.Join(city, cityLocationsFilePrefix+"-"+language+".csv"))
		if err := loader.loadCityLocationsCsv(filepath.Join(city, cityLocationsFilePrefix+"-"+language+".csv"),
			cityLocationsSql); err != nil {
			return err
		}
	}

	log.Debug().Msg("开始加载 [CountryBlocksIPv4] " + filepath.Join(country, countryIPv4BlocksFilePrefix+".csv"))
	if err := loader.loadCountryBlocksIPv4Csv(filepath.Join(country, countryIPv4BlocksFilePrefix+".csv"),
		countryBlocksIPv4Sql); err != nil {
		return err
	}

	for _, language := range Languages {
		log.Debug().Msg("开始加载 [CountryLocations-" + language + "] " + filepath.Join(country, countryLocationsFilePrefix+"-"+language+".csv"))
		if err := loader.loadCountryLocationsCsv(filepath.Join(country, countryLocationsFilePrefix+"-"+language+".csv"),
			countryLocationsSql); err != nil {
			return err
		}
	}
	log.Debug().Msg("geolite2 数据加载完成")
	return nil
}

func (loader *GeoLite2Loader) Local(asnPath, cityPath, countryPath string) error {
	return loader.loading(asnPath, cityPath, countryPath)
}

// destination 下载文件绝对路径
func (loader *GeoLite2Loader) download(url string, destination string) error {
	file, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer file.Close()

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New("下载失败，状态码：" + string(rune(response.StatusCode)))
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return errors.New("文件保存失败，错误信息：" + err.Error())
	}

	return nil
}

func (loader *GeoLite2Loader) downloader(editionID string) (string, error) {
	destination := filepath.Join(tmpDir, fmt.Sprintf("%s.%s", editionID, "zip.sha256"))
	log.Debug().Msg("开始下载 [" + editionID + "] " + fmt.Sprintf(downloadUrl, editionID, "zip.sha256"))
	if err := loader.download(fmt.Sprintf(downloadUrl, editionID, "zip.sha256"), destination); err != nil {
		return "", err
	}

	content, err := os.ReadFile(destination)
	if err != nil {
		return "", err
	}
	realHash := string(content[:64])
	filename := string(content[66 : len(content)-1])
	destination = filepath.Join(tmpDir, filename)

	log.Debug().Msg("开始下载 [" + editionID + "] " + fmt.Sprintf(downloadUrl, editionID, "zip"))
	if err := loader.download(fmt.Sprintf(downloadUrl, editionID, "zip"), destination); err != nil {
		return "", err
	}

	content, err = os.ReadFile(destination)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(content)
	hashStr := hex.EncodeToString(hash[:])

	if realHash != hashStr {
		return "", errors.New(fmt.Sprintf("sha256 不匹配，本地：%s，实际：%s", hashStr, realHash))
	}

	if err := utils.Unzip(destination, tmpDir); err != nil {
		return "", err
	}

	return filepath.Join(tmpDir, filename[:len(filename)-4]), nil
}

func (loader *GeoLite2Loader) Remote(asnEditionID, cityEditionID, countryEditionID string) error {
	asnPath, err := loader.downloader(asnEditionID)
	if err != nil {
		return err
	}
	cityPath, err := loader.downloader(cityEditionID)
	if err != nil {
		return err
	}
	countryPath, err := loader.downloader(countryEditionID)
	if err != nil {
		return err
	}

	return loader.loading(asnPath, cityPath, countryPath)
}

func (loader *GeoLite2Loader) Update() error {
	return nil
}

func (loader *GeoLite2Loader) loadASNBlocksIPv4Csv(csvPath string, sql GeoipSql) error {

	if _, err := loader.db.Exec(sql.CreateTable); err != nil {
		return err
	}

	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	tx, err := loader.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sql.Insert)
	if err != nil {
		return err
	}

	for x := 1; x < len(rows); x++ {
		start, end, err := IPRange(rows[x][0])
		if err != nil {
			return err
		}
		if _, err = stmt.Exec(rows[x][0], start, end, rows[x][1], rows[x][2]); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (loader *GeoLite2Loader) loadCityBlocksIPv4Csv(csvPath string, sql GeoipSql) error {

	if _, err := loader.db.Exec(sql.CreateTable); err != nil {
		return err
	}

	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	tx, err := loader.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sql.Insert)
	if err != nil {
		return err
	}

	for x := 1; x < len(rows); x++ {
		start, end, err := IPRange(rows[x][0])
		if err != nil {
			return err
		}
		if _, err = stmt.Exec(rows[x][0], start, end, rows[x][1], rows[x][2], rows[x][3],
			rows[x][4], rows[x][5], rows[x][6], rows[x][7], rows[x][8], rows[x][9]); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (loader *GeoLite2Loader) loadCityLocationsCsv(csvPath string, sql GeoipSql) error {

	if _, err := loader.db.Exec(sql.CreateTable); err != nil {
		return err
	}

	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	tx, err := loader.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sql.Insert)
	if err != nil {
		return err
	}

	for x := 1; x < len(rows); x++ {
		if _, err = stmt.Exec(rows[x][0], rows[x][1], rows[x][2], rows[x][3],
			rows[x][4], rows[x][5], rows[x][6], rows[x][7], rows[x][8], rows[x][9], rows[x][10], rows[x][11],
			rows[x][12], rows[x][13]); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (loader *GeoLite2Loader) loadCountryBlocksIPv4Csv(csvPath string, sql GeoipSql) error {

	if _, err := loader.db.Exec(sql.CreateTable); err != nil {
		return err
	}

	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	tx, err := loader.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sql.Insert)
	if err != nil {
		return err
	}

	for x := 1; x < len(rows); x++ {
		start, end, err := IPRange(rows[x][0])
		if err != nil {
			return err
		}
		if _, err = stmt.Exec(rows[x][0], start, end, rows[x][1], rows[x][2], rows[x][3],
			rows[x][4], rows[x][5]); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (loader *GeoLite2Loader) loadCountryLocationsCsv(csvPath string, sql GeoipSql) error {

	if _, err := loader.db.Exec(sql.CreateTable); err != nil {
		return err
	}

	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	tx, err := loader.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(sql.Insert)
	if err != nil {
		return err
	}

	for x := 1; x < len(rows); x++ {
		if _, err = stmt.Exec(rows[x][0], rows[x][1], rows[x][2], rows[x][3],
			rows[x][4], rows[x][5], rows[x][6]); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (loader *GeoLite2Loader) createDownloadRecord(sql GeoipSql) error {
	if _, err := loader.db.Exec(sql.CreateTable); err != nil {
		return err
	}
	return nil
}
