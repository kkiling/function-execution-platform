package common

import (
	"github.com/rs/zerolog/log"
	"io"
)

func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatal().Err(err).Msgf("fail close")
	}
}

func LogErr(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msgf(msg)
	}
}
