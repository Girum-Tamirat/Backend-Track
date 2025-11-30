package mongoimpl

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient holds client and database name
type MongoClient struct {
	Client *mongo.Client
	DBName string
}

func NewMongoClient(uri, dbName string) (*MongoClient, error) {
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	if dbName == "" {
		dbName = "taskdb"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	return &MongoClient{Client: client, DBName: dbName}, nil
}

func (m *MongoClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}
