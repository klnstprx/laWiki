package database

import (
	"context"
	"time"

	"github.com/laWiki/auth/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client            *mongo.Client
	UsuarioCollection *mongo.Collection
)

func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb+srv://admin:8fdCAkfmVQCRDRfK@cluster0.rfz8f.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		config.App.Logger.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		config.App.Logger.Info().Msg(config.App.MongoDBURI)
		config.App.Logger.Info().Msg(config.App.DBCollectionName)
		config.App.Logger.Fatal().Err(err).Msg("Failed to ping MongoDB")
	}

	Client = client
	UsuarioCollection = client.Database(config.App.DBName).Collection(config.App.DBCollectionName)
	config.App.Logger.Info().Msg("Connected to MongoDB")
}
