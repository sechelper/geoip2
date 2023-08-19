package geoip

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func init() {
	writer := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout})

	log.Logger = zerolog.New(writer).With().Timestamp().Caller().Logger()
}
