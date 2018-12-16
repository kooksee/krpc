package krpcs

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger zerolog.Logger

func Init() {
	logger = log.Logger.With().Str("pkg", "rpc_server").Logger()
}
