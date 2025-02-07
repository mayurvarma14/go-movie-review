package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	Client *mongo.Client
	Name   string
}

func New(ctx context.Context) (*Database, error) {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:27017/%s?authSource=%s",
		os.Getenv("MONGO_APP_USER"),
		os.Getenv("MONGO_APP_PASSWORD"),
		os.Getenv("MONGO_DOMAIN"),
		os.Getenv("MONGO_INITDB_DATABASE"),
		os.Getenv("MONGO_AUTH_SOURCE"),
	)

	client, err := mongo.Connect(options.Client().SetTimeout(10 * time.Second).ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	fmt.Println("Connected to MongoDB Successfully!")

	return &Database{
		Client: client,
		Name:   os.Getenv("MONGO_INITDB_DATABASE"),
	}, nil
}

func (db *Database) OpenCollection(name string) *mongo.Collection {
	return db.Client.Database(db.Name).Collection(name)
}
