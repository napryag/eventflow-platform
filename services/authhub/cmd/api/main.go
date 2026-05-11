package main

import (
	"fmt"
	"os"

	"github.com/napryag/eventflow-platform/libs/platform/logger"
	"github.com/napryag/eventflow-platform/services/authhub/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("failed to load config:", err)
		os.Exit(1)
	}

	log := logger.New(cfg.LogLevel)

	log.Info().Str("service", "authhub").Msg("api initialized")
}
