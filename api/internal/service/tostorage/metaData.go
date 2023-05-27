package tostorage

import (
	"github.com/kkiling/function-execution-platform/api/internal/service/model"
	"github.com/kkiling/function-execution-platform/api/internal/storage/entity"
)

func metaDataEntryToModel(v entity.MetaData) model.MetaData {
	return model.MetaData{
		Id:        v.Id,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}
