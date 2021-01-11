package mongostorage

import (
	"context"
	"go-app/server/config"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoStorage containing mongodb client connection and configuration
type MongoStorage struct {
	Config *config.DatabaseConfig
	Client *mongo.Client
}

// NewMongoStorage returns new mongodb storage instance
func NewMongoStorage(c *config.DatabaseConfig) *MongoStorage {
	clientOpts := options.Client().ApplyURI(c.ConnectionURL())
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatalf("failed to establish connection with mongodb: %s", err)
		os.Exit(1)
	}

	if err := client.Connect(context.TODO()); err != nil {
		log.Fatalf("failed to connect with mongodb: %s", err)
		os.Exit(1)
	}

	// checking if client is pining or not otherwise exit
	if err := client.Ping(context.TODO(), nil); err != nil {
		log.Fatalf("mongodb ping failed: %s", err)
		os.Exit(1)
	}
	return &MongoStorage{Config: c, Client: client}
}

// Close closes mongodb connection
func (m *MongoStorage) Close() {
	ctx := context.Background()
	m.Client.Disconnect(ctx)
}
