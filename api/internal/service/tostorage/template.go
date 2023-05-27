package tostorage

import (
	"context"
	"github.com/kkiling/function-execution-platform/api/internal/factory"
	"github.com/kkiling/function-execution-platform/api/internal/service/model"
	"github.com/kkiling/function-execution-platform/api/internal/storage/entity"
	"github.com/kkiling/id"
)

type ITemplateStorage interface {
	CountTemplate(ctx context.Context) (int64, error)
	SaveTemplate(ctx context.Context, template *model.Template) (id.Uid, error)
	UpdateTemplate(ctx context.Context, id id.Uid, template *model.Template) error
	FindTemplate(ctx context.Context, name, language string) (*model.TemplateMetaData, error)
}

func templateModelToEntry(v *model.Template) *entity.Template {
	files := make([]entity.TemplateFile, 0, len(v.Files))
	for _, f := range v.Files {
		files = append(files, entity.TemplateFile{
			FilePath: f.FilePath,
			Body:     f.Body,
		})
	}
	r := entity.Template{
		Name:               v.Name,
		Language:           v.Language,
		Version:            v.Version,
		Description:        v.Description,
		CopyFiles:          v.CopyFiles,
		ForbiddenFileNames: v.ForbiddenFileNames,
		ReadOnlyFile:       v.ReadOnlyFile,
		ContainerParams: entity.ContainerParams{
			MemoryLimitMb:       v.ContainerParams.MemoryLimitMb,
			MemoryReservationMb: v.ContainerParams.MemoryReservationMb,
			DiskSizeMb:          v.ContainerParams.DiskSizeMb,
			CPULimit:            v.ContainerParams.CPULimit,
			CPUReservation:      v.ContainerParams.CPUReservation,
			TimeoutSec:          v.ContainerParams.TimeoutSec,
		},
		Files: files,
	}
	return &r
}

func templateEntryToModel(v *entity.Template) *model.Template {
	files := make([]model.TemplateFile, 0, len(v.Files))
	for _, f := range v.Files {
		files = append(files, model.TemplateFile{
			FilePath: f.FilePath,
			Body:     f.Body,
		})
	}
	r := model.Template{
		Name:               v.Name,
		Language:           v.Language,
		Version:            v.Version,
		Description:        v.Description,
		CopyFiles:          v.CopyFiles,
		ForbiddenFileNames: v.ForbiddenFileNames,
		ReadOnlyFile:       v.ReadOnlyFile,
		ContainerParams: model.ContainerParams{
			MemoryLimitMb:       v.ContainerParams.MemoryLimitMb,
			MemoryReservationMb: v.ContainerParams.MemoryReservationMb,
			DiskSizeMb:          v.ContainerParams.DiskSizeMb,
			CPULimit:            v.ContainerParams.CPULimit,
			CPUReservation:      v.ContainerParams.CPUReservation,
			TimeoutSec:          v.ContainerParams.TimeoutSec,
		},
		Files: files,
	}
	return &r
}

type TemplateStorage struct {
	fact factory.IScopeFactory
}

func NewTemplateStorage(fact factory.IScopeFactory) *TemplateStorage {
	return &TemplateStorage{
		fact: fact,
	}
}

func (t TemplateStorage) CountTemplate(ctx context.Context) (int64, error) {
	return t.fact.TemplateStorage().CountTemplate(ctx)
}

func (t TemplateStorage) SaveTemplate(ctx context.Context, template *model.Template) (id.Uid, error) {
	return t.fact.TemplateStorage().SaveTemplate(ctx, templateModelToEntry(template))
}

func (t TemplateStorage) UpdateTemplate(ctx context.Context, id id.Uid, template *model.Template) error {
	return t.fact.TemplateStorage().UpdateTemplate(ctx, id, templateModelToEntry(template))
}

func (t TemplateStorage) FindTemplate(ctx context.Context, name, language string) (*model.TemplateMetaData, error) {
	res, err := t.fact.TemplateStorage().FindTemplate(ctx, name, language)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	return &model.TemplateMetaData{
		MetaData: metaDataEntryToModel(res.MetaData),
		Template: *templateEntryToModel(&res.Template),
	}, nil
}
