package middleware

import (
	"context"
	"time"

	"github.com/felipefrizzo/brazilian-zipcode-api/internals/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoConnection open connection with mongodb
func MongoConnection() (*mongo.Client, context.Context, error) {
	config := configs.Config.Mongo

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}
