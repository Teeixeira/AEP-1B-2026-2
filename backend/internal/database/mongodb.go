package database

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect() (*mongo.Client, error) {

	uri := os.Getenv("MONGODB_URI")

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().
		ApplyURI(uri).
		SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(opts)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	err = client.Ping(ctx, nil)

	if err != nil {
		return nil, err
	}

	return client, nil
}
