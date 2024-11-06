package database

import (
	"context"
	"time"

	"github.com/laWiki/wiki/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client         *mongo.Client
	WikiCollection *mongo.Collection
)

func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.App.MongoDBURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		config.App.Logger.Fatal().Err(err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		config.App.Logger.Fatal().Err(err)
	}

	Client = client
	WikiCollection = client.Database(config.App.DBName).Collection(config.App.DBCollectionName)
	config.App.Logger.Info().Msg("Connected to mongoDB")
}
