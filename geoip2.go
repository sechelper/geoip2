package geoip

import (
	"net"
)

type Geoip2 interface {
	//AsnBlock 查询IP ASN信息，返回 ASNBlock
	AsnBlock(ip net.IP) (*ASNBlock, error)
	//BlocksByAsnNumber 查询某组织拥有的 ASNBlock，返回 ASNBlock 数组
	BlocksByAsnNumber(int64) ([]ASNBlock, error)
	//BlocksByAsnName 查询某组织拥有的 ASNBlock，返回 ASNBlock 数组
	BlocksByAsnName(string) ([]ASNBlock, error)
	//Organizations 查询全部ASN组织，返回 Organization 数组
	Organizations() ([]Organization, error)

	//CityBlock 查询IP信息，返回 CityBlocks
	CityBlock(net.IP) (*CityBlock, error)
	//BlocksByCityCode 查询城市级某地域IP地址段，返回 CityBlock 数组
	BlocksByCityCode(language, countryCode, cityCode string) ([]CityBlock, error)

	//CountryBlock 查询IP信息，返回 CountryBlock
	CountryBlock(net.IP) (*CountryBlock, error)
	//BlocksByCountryCode 查询国家级某地域IP地址段，返回 CountryBlock 数组
	BlocksByCountryCode(string, string) ([]CountryBlock, error)
	//BlocksByContinentCode 查询洲级某地域IP地址段，返回 CountryBlock 数组
	BlocksByContinentCode(string, string) ([]CountryBlock, error)
}
