package cmd

import (
	"github.com/rs/zerolog/log"
	_ "github.com/toshke/groac/internal/executor"
)

func infoLog(msg string) {
	log.Info().Str("stage", "config").Msg(msg)
}

func groacConfig() {
	infoLog("groac: entering config stage")
}
