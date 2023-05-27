package factory

import (
	"github.com/kkiling/function-execution-platform/api/internal/config"
	"github.com/kkiling/function-execution-platform/api/internal/service"
	"github.com/kkiling/function-execution-platform/api/internal/storage"
	"github.com/kkiling/id"
)

type singletonFactoryBase interface {
	GetConfig() config.Config
	TemplateStorage() storage.ITemplateStorage
	GeneratorId() id.IGeneratorId
}

type ISingletonFactory interface {
	singletonFactoryBase
	CreateScopeFactory() IScopeFactory
}

type IScopeFactory interface {
	singletonFactoryBase
	GetTemplateService() service.ITemplateService
}
