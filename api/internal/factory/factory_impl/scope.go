package factory_impl

import (
	"github.com/kkiling/function-execution-platform/api/internal/config"
	"github.com/kkiling/function-execution-platform/api/internal/factory"
	"github.com/kkiling/function-execution-platform/api/internal/service"
	"github.com/kkiling/function-execution-platform/api/internal/service/template_impl"
	"github.com/kkiling/function-execution-platform/api/internal/storage"
	"github.com/kkiling/id"
)

type ScopeFactory struct {
	fact     factory.ISingletonFactory
	template service.ITemplateService
}

func NewScopeFactory(fact factory.ISingletonFactory) *ScopeFactory {
	res := ScopeFactory{
		fact: fact,
	}
	return &res
}

func (f *ScopeFactory) GetConfig() config.Config {
	return f.fact.GetConfig()
}

func (f *ScopeFactory) TemplateStorage() storage.ITemplateStorage {
	return f.fact.TemplateStorage()
}

func (f *ScopeFactory) GeneratorId() id.IGeneratorId {
	return f.fact.GeneratorId()
}

// *** SCOPE

func (f *ScopeFactory) GetTemplateService() service.ITemplateService {
	if f.template == nil {
		f.template = template_impl.NewService(f)
	}
	return f.template
}
