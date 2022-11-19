package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"

	v1 "refactoring/internal/api/v1"
	"refactoring/internal/config"
	"refactoring/internal/storage"
)

func main() {

	newCfg, err := config.New(config.WithFlag, config.WithEnv)
	if err != nil {
		log.Error().Err(err).Msg("creating config")
		os.Exit(1)
	}

	newStorage, err := storage.New(newCfg)
	if err != nil {
		log.Error().Err(err).Msg("creating storage")
		os.Exit(1)
	}

	newAPI, err := v1.New(newCfg, newStorage)
	if err != nil {
		log.Error().Err(err).Msg("creating api")
		os.Exit(1)
	}

	ctx := context.Background()

	err = newAPI.Run(ctx)
	if err != nil {
		log.Error().Err(err).Msg("running api")
		os.Exit(1)
	}

	os.Exit(0)
}
