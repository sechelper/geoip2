package geoip

import (
	"database/sql"
	"net"
)

type Geolite2 struct {
	db *sql.DB
}

func NewGeolite2(db *sql.DB) Geoip2 {
	return Geolite2{db: db}
}

func (geo Geolite2) AsnBlock(ip net.IP) (*ASNBlock, error) {
	row := geo.db.QueryRow("SELECT network, autonomous_system_number, autonomous_system_organization "+
		"FROM GeoLite2ASNBlocksIPv4 WHERE ? BETWEEN start_ip AND end_ip", IP2Int(ip))

	var blocks = new(ASNBlock)
	if err := row.Scan(&blocks.Network, &blocks.AutonomousSystemNumber,
		&blocks.AutonomousSystemOrganization); err != nil {
		return nil, err
	}

	return blocks, nil
}

func (geo Geolite2) BlocksByAsnNumber(number int64) ([]ASNBlock, error) {
	rows, err := geo.db.Query("SELECT network, autonomous_system_number, autonomous_system_organization "+
		"FROM GeoLite2ASNBlocksIPv4 WHERE autonomous_system_number=?", number)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []ASNBlock
	for rows.Next() {
		var block ASNBlock
		if err := rows.Scan(&block.Network, &block.AutonomousSystemNumber,
			&block.AutonomousSystemOrganization); err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return blocks, nil
}

func (geo Geolite2) BlocksByAsnName(name string) ([]ASNBlock, error) {
	rows, err := geo.db.Query("SELECT network, autonomous_system_number, autonomous_system_organization "+
		"FROM GeoLite2ASNBlocksIPv4 WHERE autonomous_system_organization=?", name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []ASNBlock

	for rows.Next() {
		var block ASNBlock
		if err := rows.Scan(&block.Network, &block.AutonomousSystemNumber,
			&block.AutonomousSystemOrganization); err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return blocks, nil
}

func (geo Geolite2) Organizations() ([]Organization, error) {
	rows, err := geo.db.Query("SELECT autonomous_system_number, autonomous_system_organization " +
		"FROM GeoLite2ASNBlocksIPv4 GROUP BY autonomous_system_number;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orgs []Organization

	for rows.Next() {
		var org Organization
		if err := rows.Scan(&org.AutonomousSystemNumber,
			&org.AutonomousSystemOrganization); err != nil {
			return nil, err
		}
		orgs = append(orgs, org)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orgs, nil
}

func (geo Geolite2) CityBlock(ip net.IP) (*CityBlock, error) {
	row := geo.db.QueryRow("SELECT GeoLite2CityBlocksIPv4.network,GeoLite2CityBlocksIPv4.geoname_id,GeoLite2CityBlocksIPv4.registered_country_geoname_id,GeoLite2CityBlocksIPv4."+
		"represented_country_geoname_id,GeoLite2CityBlocksIPv4.is_anonymous_proxy,GeoLite2CityBlocksIPv4.is_satellite_provider,GeoLite2CityBlocksIPv4.postal_code,GeoLite2CityBlocksIPv4.latitude,GeoLite2CityBlocksIPv4.longitude,"+
		"GeoLite2CityBlocksIPv4.accuracy_radius,GeoLite2CityLocations.geoname_id,GeoLite2CityLocations.locale_code,GeoLite2CityLocations.continent_code,GeoLite2CityLocations.continent_name,"+
		"GeoLite2CityLocations.country_iso_code,GeoLite2CityLocations.country_name,GeoLite2CityLocations.subdivision_1_iso_code,GeoLite2CityLocations.subdivision_1_name,"+
		"GeoLite2CityLocations.subdivision_2_iso_code,GeoLite2CityLocations.subdivision_2_name,GeoLite2CityLocations.city_name,GeoLite2CityLocations.metro_code,GeoLite2CityLocations.time_zone,"+
		"GeoLite2CityLocations.is_in_european_union FROM GeoLite2CityBlocksIPv4 "+
		"LEFT JOIN GeoLite2CityLocations ON GeoLite2CityBlocksIPv4.geoname_id = GeoLite2CityLocations.geoname_id"+
		" WHERE ? BETWEEN start_ip AND end_ip ",
		IP2Int(ip))

	var block = new(CityBlock)
	block.location = new(CityLocation)
	if err := row.Scan(&block.Network, &block.GeonameID, &block.RegisteredCountryGeonameID,
		&block.RepresentedCountryGeonameID, &block.IsAnonymousProxy, &block.IsSatelliteProvider, &block.PostalCode,
		&block.Latitude, &block.Longitude, &block.AccuracyRadius, &block.location.GeonameID, &block.location.LocaleCode,
		&block.location.ContinentCode, &block.location.ContinentName, &block.location.CountryISOCode,
		&block.location.CountryName, &block.location.Subdivision1ISOCode, &block.location.Subdivision1Name,
		&block.location.Subdivision2ISOCode, &block.location.Subdivision2Name, &block.location.CityName,
		&block.location.MetroCode, &block.location.TimeZone, &block.location.IsInEuropeanUnion); err != nil {
		return nil, err
	}

	return block, nil
}

func (geo Geolite2) BlocksByCityCode(language, countryCode, cityCode string) ([]CityBlock, error) {
	rows, err := geo.db.Query("SELECT GeoLite2CityBlocksIPv4.network,GeoLite2CityBlocksIPv4.geoname_id,GeoLite2CityBlocksIPv4.registered_country_geoname_id,GeoLite2CityBlocksIPv4."+
		"represented_country_geoname_id,GeoLite2CityBlocksIPv4.is_anonymous_proxy,GeoLite2CityBlocksIPv4.is_satellite_provider,GeoLite2CityBlocksIPv4.postal_code,GeoLite2CityBlocksIPv4.latitude,GeoLite2CityBlocksIPv4.longitude,"+
		"GeoLite2CityBlocksIPv4.accuracy_radius,GeoLite2CityLocations.geoname_id,GeoLite2CityLocations.locale_code,GeoLite2CityLocations.continent_code,GeoLite2CityLocations.continent_name,"+
		"GeoLite2CityLocations.country_iso_code,GeoLite2CityLocations.country_name,GeoLite2CityLocations.subdivision_1_iso_code,GeoLite2CityLocations.subdivision_1_name,"+
		"GeoLite2CityLocations.subdivision_2_iso_code,GeoLite2CityLocations.subdivision_2_name,GeoLite2CityLocations.city_name,GeoLite2CityLocations.metro_code,GeoLite2CityLocations.time_zone,"+
		"GeoLite2CityLocations.is_in_european_union FROM GeoLite2CityBlocksIPv4 "+
		"LEFT JOIN GeoLite2CityLocations ON GeoLite2CityBlocksIPv4.geoname_id = GeoLite2CityLocations.geoname_id"+
		" WHERE locale_code=? and country_iso_code=? and subdivision_1_iso_code=?", language, countryCode, cityCode)
	if err != nil { //
		return nil, err
	}
	defer rows.Close()

	var blocks []CityBlock

	for rows.Next() {
		var block CityBlock
		block.location = new(CityLocation)
		if err := rows.Scan(&block.Network, &block.GeonameID, &block.RegisteredCountryGeonameID,
			&block.RepresentedCountryGeonameID, &block.IsAnonymousProxy, &block.IsSatelliteProvider, &block.PostalCode,
			&block.Latitude, &block.Longitude, &block.AccuracyRadius, &block.location.GeonameID, &block.location.LocaleCode,
			&block.location.ContinentCode, &block.location.ContinentName, &block.location.CountryISOCode,
			&block.location.CountryName, &block.location.Subdivision1ISOCode, &block.location.Subdivision1Name,
			&block.location.Subdivision2ISOCode, &block.location.Subdivision2Name, &block.location.CityName,
			&block.location.MetroCode, &block.location.TimeZone, &block.location.IsInEuropeanUnion); err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return blocks, nil
}

func (geo Geolite2) CountryBlock(ip net.IP) (*CountryBlock, error) {
	row := geo.db.QueryRow("SELECT GeoLite2CountryBlocksIPv4.network,"+
		"GeoLite2CountryBlocksIPv4.geoname_id,GeoLite2CountryBlocksIPv4.registered_country_geoname_id,"+
		"GeoLite2CountryBlocksIPv4.represented_country_geoname_id,GeoLite2CountryBlocksIPv4.is_anonymous_proxy,"+
		"GeoLite2CountryBlocksIPv4.is_satellite_provider,"+
		"GeoLite2CountryLocations.geoname_id,GeoLite2CountryLocations.locale_code,GeoLite2CountryLocations.continent_code,"+
		"GeoLite2CountryLocations.continent_name,GeoLite2CountryLocations.country_iso_code,"+
		"GeoLite2CountryLocations.country_name,GeoLite2CountryLocations.is_in_european_union "+
		"FROM GeoLite2CountryBlocksIPv4 "+
		"LEFT JOIN GeoLite2CountryLocations ON GeoLite2CountryBlocksIPv4.geoname_id = GeoLite2CountryLocations.geoname_id "+
		"WHERE  ? BETWEEN start_ip AND end_ip",
		IP2Int(ip))

	var block = new(CountryBlock)
	block.location = new(CountryLocation)
	if err := row.Scan(&block.Network, &block.GeonameID, &block.RegisteredCountryGeonameID,
		&block.RepresentedCountryGeonameID, &block.IsAnonymousProxy, &block.IsSatelliteProvider,
		&block.location.GeonameID, &block.location.LocaleCode, &block.location.ContinentCode,
		&block.location.ContinentName, &block.location.CountryISOCode, &block.location.ContinentName,
		&block.location.IsInEuropeanUnion); err != nil {
		return nil, err
	}

	return block, nil
}

func (geo Geolite2) BlocksByCountryCode(language, code string) ([]CountryBlock, error) {
	rows, err := geo.db.Query("SELECT GeoLite2CountryBlocksIPv4.network,"+
		"GeoLite2CountryBlocksIPv4.geoname_id,GeoLite2CountryBlocksIPv4.registered_country_geoname_id,"+
		"GeoLite2CountryBlocksIPv4.represented_country_geoname_id,GeoLite2CountryBlocksIPv4.is_anonymous_proxy,"+
		"GeoLite2CountryBlocksIPv4.is_satellite_provider,"+
		"GeoLite2CountryLocations.geoname_id,GeoLite2CountryLocations.locale_code,GeoLite2CountryLocations.continent_code,"+
		"GeoLite2CountryLocations.continent_name,GeoLite2CountryLocations.country_iso_code,"+
		"GeoLite2CountryLocations.country_name,GeoLite2CountryLocations.is_in_european_union "+
		"FROM GeoLite2CountryBlocksIPv4 "+
		"LEFT JOIN GeoLite2CountryLocations ON GeoLite2CountryBlocksIPv4.geoname_id = GeoLite2CountryLocations.geoname_id "+
		"WHERE GeoLite2CountryLocations.locale_code=? and country_iso_code=?", language, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []CountryBlock

	for rows.Next() {
		var block = new(CountryBlock)
		block.location = new(CountryLocation)
		if err := rows.Scan(&block.Network, &block.GeonameID, &block.RegisteredCountryGeonameID,
			&block.RepresentedCountryGeonameID, &block.IsAnonymousProxy, &block.IsSatelliteProvider,
			&block.location.GeonameID, &block.location.LocaleCode, &block.location.ContinentCode,
			&block.location.ContinentName, &block.location.CountryISOCode, &block.location.ContinentName,
			&block.location.IsInEuropeanUnion); err != nil {
			return nil, err
		}
		blocks = append(blocks, *block)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return blocks, nil
}

func (geo Geolite2) BlocksByContinentCode(language, code string) ([]CountryBlock, error) {
	rows, err := geo.db.Query("SELECT GeoLite2CountryBlocksIPv4.network,"+
		"GeoLite2CountryBlocksIPv4.geoname_id,GeoLite2CountryBlocksIPv4.registered_country_geoname_id,"+
		"GeoLite2CountryBlocksIPv4.represented_country_geoname_id,GeoLite2CountryBlocksIPv4.is_anonymous_proxy,"+
		"GeoLite2CountryBlocksIPv4.is_satellite_provider,"+
		"GeoLite2CountryLocations.geoname_id,GeoLite2CountryLocations.locale_code,GeoLite2CountryLocations.continent_code,"+
		"GeoLite2CountryLocations.continent_name,GeoLite2CountryLocations.country_iso_code,"+
		"GeoLite2CountryLocations.country_name,GeoLite2CountryLocations.is_in_european_union "+
		"FROM GeoLite2CountryBlocksIPv4 "+
		"LEFT JOIN GeoLite2CountryLocations ON GeoLite2CountryBlocksIPv4.geoname_id = GeoLite2CountryLocations.geoname_id "+
		"WHERE GeoLite2CountryLocations.locale_code=? and continent_code=?", language, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []CountryBlock

	for rows.Next() {
		var block = new(CountryBlock)
		block.location = new(CountryLocation)
		if err := rows.Scan(&block.Network, &block.GeonameID, &block.RegisteredCountryGeonameID,
			&block.RepresentedCountryGeonameID, &block.IsAnonymousProxy, &block.IsSatelliteProvider,
			&block.location.GeonameID, &block.location.LocaleCode, &block.location.ContinentCode,
			&block.location.ContinentName, &block.location.CountryISOCode, &block.location.ContinentName,
			&block.location.IsInEuropeanUnion); err != nil {
			return nil, err
		}
		blocks = append(blocks, *block)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return blocks, nil
}
