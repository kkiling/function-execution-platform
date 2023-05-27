package entity

import (
	"github.com/kkiling/id"
	"time"
)

type MetaData struct {
	Id        id.Uid    `bson:"_id"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}
