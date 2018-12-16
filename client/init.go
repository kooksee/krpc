package krpcc

import (
	logger "github.com/rs/zerolog/log"
	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init() {
	log = logger.Logger.With().Str("pkg", "rpcclient").Logger()
}

type M map[string]interface{}

type P struct {
	Method string
	Params map[string]interface{}
	Result interface{}
}

func Call() {

}
