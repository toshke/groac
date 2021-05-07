package cmd

import (
	"github.com/rs/zerolog/log"
)

func groacCleanup() {
	log.Info().Str("stage", "cleanup").Msg("groac: entering cleanup stage")
}
