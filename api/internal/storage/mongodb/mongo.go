package mongodb

import (
	"context"
	"github.com/kkiling/id"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	templateCollections string = "template"
	maxPoolSize         uint64 = 100
	minPoolSize         uint64 = 10
)

type Cfg struct {
	Uri      string
	Database string
}

type MongodbStorage struct {
	cancel      context.CancelFunc
	generatorId id.IGeneratorId
	client      *mongo.Client
	db          *mongo.Database
	cfg         Cfg
}

func NewMongodbStorage(ctx context.Context, cfg Cfg, generatorId id.IGeneratorId) (*MongodbStorage, error) {
	clientOptions := options.Client().ApplyURI(cfg.Uri)
	clientOptions.SetMaxPoolSize(maxPoolSize)
	clientOptions.SetMinPoolSize(minPoolSize)
	clientOptions.SetMaxConnecting(minPoolSize)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, "fail connect to mongodb")
	}

	return &MongodbStorage{
		client:      client,
		db:          client.Database(cfg.Database),
		cfg:         cfg,
		generatorId: generatorId,
	}, nil
}

func (d *MongodbStorage) createTemplateIndex(ctx context.Context, collectionNames []string) error {
	collectionExist := false
	for _, name := range collectionNames {
		if name == templateCollections {
			collectionExist = true
			break
		}
	}

	if collectionExist {
		return nil
	}

	index := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}, {Key: "language", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	opts := options.CreateIndexes()
	_, err := d.db.Collection(templateCollections).Indexes().CreateMany(ctx, index, opts)
	if err != nil {
		return errors.Wrap(err, "fail create template index")
	}

	return nil
}

func (d *MongodbStorage) Shutdown(ctx context.Context) error {
	err := d.client.Disconnect(ctx)
	if err != nil {
		return errors.Wrap(err, "fail disconnect mongodb")
	}
	return nil
}
