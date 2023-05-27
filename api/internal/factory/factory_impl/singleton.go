package factory_impl

import (
	"context"
	"github.com/kkiling/function-execution-platform/api/internal/config"
	"github.com/kkiling/function-execution-platform/api/internal/factory"
	"github.com/kkiling/function-execution-platform/api/internal/storage"
	"github.com/kkiling/function-execution-platform/api/internal/storage/mongodb"
	"github.com/kkiling/id"
	generator_id_impl "github.com/kkiling/id/id_impl"
	"github.com/pkg/errors"
)

type SingletonFactory struct {
	cfg         *config.Config
	generatorId id.IGeneratorId
	mongo       *mongodb.MongodbStorage
}

func NewSingletonFactory(ctx context.Context, configFile string) (*SingletonFactory, error) {
	// Load config
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		return nil, errors.Wrap(err, "fail to initialize config")
	}

	generatorId, err := generator_id_impl.NewService(cfg.MachineID)
	if err != nil {
		return nil, errors.Wrap(err, "fail to initialize generator id")
	}

	mongo, err := mongodb.NewMongodbStorage(ctx, mongodb.Cfg{
		Uri:      cfg.Mongo.Uri,
		Database: cfg.Mongo.Database,
	}, generatorId)
	if err != nil {
		return nil, errors.Wrap(err, "fail create mongo")
	}

	err = mongo.CreateIndexes(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "fail to create indexes in mongodb")
	}

	return &SingletonFactory{
		cfg:         cfg,
		generatorId: generatorId,
		mongo:       mongo,
	}, nil
}

func (f *SingletonFactory) CreateScopeFactory() factory.IScopeFactory {
	return NewScopeFactory(f)
}

func (f *SingletonFactory) GetConfig() config.Config {
	return *f.cfg
}

func (f *SingletonFactory) TemplateStorage() storage.ITemplateStorage {
	return f.mongo
}

func (f *SingletonFactory) GeneratorId() id.IGeneratorId {
	return f.generatorId
}

func (f *SingletonFactory) Shutdown(ctx context.Context) error {
	if err := f.mongo.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "fail shutdown mongodb")
	}
	return nil
}
