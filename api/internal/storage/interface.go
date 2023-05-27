package storage

import (
	"context"
	"github.com/kkiling/function-execution-platform/api/internal/storage/entity"
	"github.com/kkiling/id"
)

type ITemplateStorage interface {
	CountTemplate(ctx context.Context) (int64, error)
	SaveTemplate(ctx context.Context, template *entity.Template) (id.Uid, error)
	UpdateTemplate(ctx context.Context, id id.Uid, template *entity.Template) error
	FindTemplate(ctx context.Context, name, language string) (*entity.TemplateMetaData, error)
}
