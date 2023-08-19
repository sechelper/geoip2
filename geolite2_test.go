package geoip

import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog"
	"log"
	"net"
	"os"
	"testing"
)

var loader *GeoLite2Loader
var geo Geoip2

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// Open a database connection
	db, err := sql.Open("sqlite3", "geoip2-test.db")
	if err != nil {
		log.Fatal(err)
	}

	loader = NewGeoLite2Loader(db)

	_, err = os.Stat("geoip2-test.db")

	geo = NewGeolite2(db)
}
func TestGeoLite2Loader_Local(t *testing.T) {
	//if err := loader.Local("/Users/kun/Downloads/GeoLite2-ASN-CSV_20230815",
	//	"/Users/kun/Downloads/GeoLite2-City-CSV_20230815",
	//	"/Users/kun/Downloads/GeoLite2-Country-CSV_20230815"); err != nil {
	//	log.Fatal(err)
	//}
}

func TestGeoLite2Loader_Remote(t *testing.T) {
	//err := loader.Remote("GeoLite2-ASN-CSV", "GeoLite2-City-CSV", "GeoLite2-Country-CSV")
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func TestGeolite2_AsnBlock(t *testing.T) {
	block, err := geo.AsnBlock(net.ParseIP("59.110.190.34"))
	if err != nil {
		return
	}
	fmt.Println(block)
}

func TestGeolite2_BlocksByAsnNumber(t *testing.T) {
	blocks, err := geo.BlocksByAsnNumber(38803)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(blocks)
}

func TestGeolite2_BlocksByAsnName(t *testing.T) {
	blocks, err := geo.BlocksByAsnName("Wirefreebroadband Pty Ltd")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(blocks)
}

func TestGeolite2_Organizations(t *testing.T) {
	orgs, err := geo.Organizations()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(orgs)
}

func TestGeolite2_CityBlock(t *testing.T) {
	block, err := geo.CityBlock(net.ParseIP("59.110.190.34"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*block.location)
}

func TestGeolite2_CountryBlock(t *testing.T) {
	block, err := geo.CountryBlock(net.ParseIP("59.110.190.34"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*block.location)
}

func TestGeolite2_BlocksByCityCode(t *testing.T) {
	blocks, err := geo.BlocksByCityCode("zh-CN", "CN", "BJ")
	if err != nil {
		log.Fatal(err)
	}
	for _, blocks := range blocks {
		fmt.Println(*blocks.location)
	}
}
func TestGeolite2_BlocksByCountryCode(t *testing.T) {
	blocks, err := geo.BlocksByCountryCode("zh-CN", "CN")
	if err != nil {
		log.Fatal(err)
	}
	for _, block := range blocks {
		fmt.Println(block)
	}
}

func TestGeolite2_BlocksByContinentCode(t *testing.T) {
	blocks, err := geo.BlocksByCountryCode("zh-CN", "AS")
	if err != nil {
		log.Fatal(err)
	}
	for _, block := range blocks {
		fmt.Println(block)
	}
}
