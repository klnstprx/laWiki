package database

import (
	"context"
	"os"
	"time"

	"github.com/laWiki/comment/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client            *mongo.Client
	CommentCollection *mongo.Collection
)

func Connect() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
		config.App.Logger.Warn().Str("mongoURI", mongoURI).Msg("MongoURI not set in .env, defaulting to locally hosted DB.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
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
	CommentCollection = client.Database("laWiki").Collection("comments")
	config.App.Logger.Info().Msg("Connected to mongoDB")
}
