package krpcc

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger zerolog.Logger

func Init() {
	logger = log.Logger.With().Str("pkg", "rpc_client").Logger()
}

type M map[string]interface{}

type P struct {
	Method string
	Params map[string]interface{}
	Result interface{}
}
