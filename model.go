package geoip

type ASNBlock struct {
	Network string `json:"network"`
	Organization
}

type Organization struct {
	AutonomousSystemNumber       int    `json:"autonomous_system_number"`
	AutonomousSystemOrganization string `json:"autonomous_system_organization"`
}

type CityBlock struct {
	ID                          int64   `json:"id"`
	Network                     string  `json:"network"`
	GeonameID                   int64   `json:"geoname_id"`
	RegisteredCountryGeonameID  string  `json:"registered_country_geoname_id"`
	RepresentedCountryGeonameID string  `json:"represented_country_geoname_id"`
	IsAnonymousProxy            int     `json:"is_anonymous_proxy"`
	IsSatelliteProvider         int     `json:"is_satellite_provider"`
	PostalCode                  string  `json:"postal_code"`
	Latitude                    float64 `json:"latitude"`
	Longitude                   float64 `json:"longitude"`
	AccuracyRadius              int     `json:"accuracy_radius"`
	location                    *CityLocation
}

type CityLocation struct {
	ID                  int64  `json:"id"`
	GeonameID           int64  `json:"geoname_id"`
	LocaleCode          string `json:"locale_code"`
	ContinentCode       string `json:"continent_code"`
	ContinentName       string `json:"continent_name"`
	CountryISOCode      string `json:"country_iso_code"`
	CountryName         string `json:"country_name"`
	Subdivision1ISOCode string `json:"subdivision_1_iso_code"`
	Subdivision1Name    string `json:"subdivision_1_name"`
	Subdivision2ISOCode string `json:"subdivision_2_iso_code"`
	Subdivision2Name    string `json:"subdivision_2_name"`
	CityName            string `json:"city_name"`
	MetroCode           string `json:"metro_code"`
	TimeZone            string `json:"time_zone"`
	IsInEuropeanUnion   string `json:"is_in_european_union"`
}

type CountryBlock struct {
	ID                          int64  `json:"id"`
	Network                     string `json:"network"`
	GeonameID                   int64  `json:"geoname_id"`
	RegisteredCountryGeonameID  string `json:"registered_country_geoname_id"`
	RepresentedCountryGeonameID string `json:"represented_country_geoname_id"`
	IsAnonymousProxy            string `json:"is_anonymous_proxy"`
	IsSatelliteProvider         string `json:"is_satellite_provider"`
	location                    *CountryLocation
}

type CountryLocation struct {
	ID                int64  `json:"id"`
	GeonameID         int64  `json:"geoname_id"`
	LocaleCode        string `json:"locale_code"`
	ContinentCode     string `json:"continent_code"`
	ContinentName     string `json:"continent_name"`
	CountryISOCode    string `json:"country_iso_code"`
	CountryName       string `json:"country_name"`
	IsInEuropeanUnion string `json:"is_in_european_union"`
}
