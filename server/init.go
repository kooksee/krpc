package krpcs

import (
	logger "github.com/rs/zerolog/log"
	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Init() {
	log = logger.Logger.With().Str("pkg", "rpcserver").Logger()
}
