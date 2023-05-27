package mongodb

import (
	"context"
	"github.com/kkiling/function-execution-platform/api/internal/storage/entity"
	"github.com/kkiling/id"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (d *MongodbStorage) CreateIndexes(ctx context.Context) error {
	names, err := d.db.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return errors.Wrap(err, "fail get list collections")
	}

	err = d.createTemplateIndex(ctx, names)
	if err != nil {
		return errors.Wrap(err, "fail create template index")
	}
	return nil
}

func (d *MongodbStorage) CountTemplate(ctx context.Context) (int64, error) {
	count, err := d.db.Collection(templateCollections).CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (d *MongodbStorage) FindTemplate(ctx context.Context, name, language string) (*entity.TemplateMetaData, error) {
	var template entity.TemplateMetaData
	filter := bson.D{{"name", name}, {"language", language}}

	err := d.db.Collection(templateCollections).FindOne(ctx, filter).Decode(&template)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &template, nil
}

func (d *MongodbStorage) UpdateTemplate(ctx context.Context, id id.Uid, template *entity.Template) error {
	filter := bson.D{{"_id", id}}

	decodedTemplate, err := objectToBsonMap(template)
	if err != nil {
		return err
	}
	decodedTemplate["updatedAt"] = time.Now()

	update := bson.D{{"$set", decodedTemplate}}
	_, err = d.db.Collection(templateCollections).UpdateOne(ctx, filter, update)
	return err
}

func (d *MongodbStorage) SaveTemplate(ctx context.Context, template *entity.Template) (id.Uid, error) {
	now := time.Now()
	t := entity.TemplateMetaData{
		MetaData: entity.MetaData{
			Id:        d.generatorId.NextId(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Template: *template,
	}

	_, err := d.db.Collection(templateCollections).InsertOne(ctx, t)
	return t.Id, err
}
