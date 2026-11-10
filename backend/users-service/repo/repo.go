package repo

import (
	"2x4km/users/config"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type UsersRepository struct {
	*mongo.Collection
}

func InitDB(ctx context.Context, cfg config.DBConfig) (*UsersRepository, func(ctx context.Context) error, error) {
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

	return &UsersRepository{client.Database("users_db").Collection("users")}, client.Disconnect, nil
}
