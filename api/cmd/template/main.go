package main

import (
	"context"
	"flag"
	"github.com/kkiling/function-execution-platform/api/internal/config"
	"github.com/kkiling/function-execution-platform/api/internal/factory/factory_impl"
	"github.com/kkiling/function-execution-platform/api/pkg/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	err := logging.InitLogging(&logging.Config{
		SystemName: config.Namespace,
		Level:      zerolog.InfoLevel,
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	configFile := flag.String("config", "configs/config.yml", "Path to config file")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fact, err := factory_impl.NewSingletonFactory(ctx, *configFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("fail create factory")
	}

	err = fact.CreateScopeFactory().GetTemplateService().InitBaseTemplates(ctx)
	if err != nil {
		log.Fatal().Err(err).Msgf("fail init base template")
	}

	err = fact.CreateScopeFactory().GetTemplateService().LoadGitTemplates(ctx, "git@github.com:kkiling/function-template.git", "main")
	if err != nil {
		log.Fatal().Err(err).Msgf("fail init base template")
	}
}
