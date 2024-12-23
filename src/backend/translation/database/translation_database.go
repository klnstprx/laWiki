package database

import (
	"context"
	"time"

	"github.com/laWiki/translation/config"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client                *mongo.Client
	TranslationCollection *mongo.Collection
)

// Connect initializes the MongoDB connection
func Connect(cfg config.GlobalConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDBURI)
	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	// Verificar la conexión
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ping MongoDB")
	}

	TranslationCollection = Client.Database(cfg.DBName).Collection("translations")
	log.Info().Msg("Connected to MongoDB")
}
