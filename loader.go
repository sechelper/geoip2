package geoip

import (
	_ "github.com/mattn/go-sqlite3"
)

type GeoipSql struct {
	CreateTable string
	Insert      string
}
