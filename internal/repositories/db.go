package repositories

import (
	"context"
	"github.com/EtienneBerube/cat-scribers/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	connectTimeout = 5
)

// getDBConnection returns a connection to the MongoDB instance provided by the MongoDBURL
func getDBConnection() (*mongo.Client, context.Context, context.CancelFunc) {

	client, err := mongo.NewClient(options.Client().ApplyURI(config.Config.MongoDBURL))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to cluster: %v", err)
	}

	return client, ctx, cancel
}
