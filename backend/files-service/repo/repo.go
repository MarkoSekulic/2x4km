package repo

import (
	"2x4km/files/config"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const DATABASE_NAME = "files_db"

type FilesRepository struct {
	client *mongo.Client
	db     *mongo.Database
}

// Client vraca mongo klijenta (potreban za StartSession / transakcije).
func (r *FilesRepository) Client() *mongo.Client {
	return r.client
}

// Collection vraca kolekciju iz baze files_db po imenu (root, folders, files).
func (r *FilesRepository) Collection(name string) *mongo.Collection {
	return r.db.Collection(name)
}

func InitDB(ctx context.Context, cfg config.DBConfig) (*FilesRepository, func(ctx context.Context) error, error) {
	clientOptions := options.Client().ApplyURI(cfg.Address).SetMinPoolSize(cfg.MinPool).SetMaxPoolSize(cfg.MaxPool)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to mongo: %w", err)
	}

	withTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	err = client.Ping(withTimeout, readpref.Primary())
	if err != nil {
		client.Disconnect(ctx)
		return nil, nil, fmt.Errorf("ping mongo: %w", err)
	}

	return &FilesRepository{client: client, db: client.Database(DATABASE_NAME)}, client.Disconnect, nil
}
