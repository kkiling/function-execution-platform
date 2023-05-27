package model

import (
	"github.com/kkiling/id"
	"time"
)

type MetaData struct {
	Id        id.Uid
	CreatedAt time.Time
	UpdatedAt time.Time
}
