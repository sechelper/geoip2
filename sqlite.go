package geoip

var asnBlocksIPv4Sql = GeoipSql{
	CreateTable: `CREATE TABLE GeoLite2ASNBlocksIPv4 (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    network TEXT,
    start_ip INTEGER,
    end_ip INTEGER,
    autonomous_system_number TEXT,
    autonomous_system_organization TEXT
);`,
	Insert: `INSERT INTO GeoLite2ASNBlocksIPv4 (network, start_ip, end_ip, autonomous_system_number, autonomous_system_organization) VALUES(?, ?, ?, ?, ?)`,
}

var cityBlocksIPv4Sql = GeoipSql{
	CreateTable: `CREATE TABLE GeoLite2CityBlocksIPv4 (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      network TEXT,
      start_ip INTEGER,
      end_ip INTEGER,
      geoname_id INTEGER,
      registered_country_geoname_id TEXT,
      represented_country_geoname_id TEXT,
      is_anonymous_proxy INTEGER,
      is_satellite_provider INTEGER,
      postal_code TEXT,
      latitude REAL,
      longitude REAL,
      accuracy_radius INTEGER
);`,
	Insert: `INSERT INTO GeoLite2CityBlocksIPv4 (network, start_ip, end_ip, geoname_id, registered_country_geoname_id, represented_country_geoname_id, 
            is_anonymous_proxy, is_satellite_provider, postal_code, latitude, longitude, accuracy_radius)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`,
}

var cityLocationsSql = GeoipSql{
	CreateTable: `CREATE TABLE IF NOT EXISTS GeoLite2CityLocations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    geoname_id INTEGER,
    locale_code TEXT,
    continent_code TEXT,
    continent_name TEXT,
    country_iso_code TEXT,
    country_name TEXT,
    subdivision_1_iso_code TEXT,
    subdivision_1_name TEXT,
    subdivision_2_iso_code TEXT,
    subdivision_2_name TEXT,
    city_name TEXT,
    metro_code TEXT,
    time_zone TEXT,
    is_in_european_union TEXT,
    language TEXT                    
);`,
	Insert: `INSERT INTO GeoLite2CityLocations(geoname_id, locale_code, continent_code, continent_name, country_iso_code, country_name, subdivision_1_iso_code, subdivision_1_name, subdivision_2_iso_code, subdivision_2_name, city_name, metro_code, time_zone, is_in_european_union) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
}

var countryBlocksIPv4Sql = GeoipSql{
	CreateTable: `CREATE TABLE GeoLite2CountryBlocksIPv4 (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    network TEXT,
    start_ip INTEGER,
    end_ip INTEGER,
    geoname_id INTEGER,
    registered_country_geoname_id TEXT,
    represented_country_geoname_id TEXT,
    is_anonymous_proxy TEXT,
    is_satellite_provider TEXT
);`,
	Insert: `INSERT INTO GeoLite2CountryBlocksIPv4 (network, start_ip, end_ip, geoname_id, registered_country_geoname_id, represented_country_geoname_id, is_anonymous_proxy, is_satellite_provider) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`,
}

var countryLocationsSql = GeoipSql{
	CreateTable: `CREATE TABLE IF NOT EXISTS GeoLite2CountryLocations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  geoname_id INTEGER,
  locale_code TEXT,
  continent_code TEXT,
  continent_name TEXT,
  country_iso_code TEXT,
  country_name TEXT,
  is_in_european_union INTEGER,
  language TEXT
);`,
	Insert: `INSERT INTO GeoLite2CountryLocations (geoname_id, locale_code, continent_code, continent_name, country_iso_code, country_name, is_in_european_union) VALUES (?, ?, ?, ?, ?, ?, ?)`,
}
