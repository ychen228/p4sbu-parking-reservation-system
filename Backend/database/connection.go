package mongodb

import (
	"context"
	"fmt"
	"log"

	"github.com/YimingChen/P4SBU/Backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connect establishes a connection to MongoDB
func Connect(cfg *config.Config) (*mongo.Client, error) {
	var connectionString string

	// If credentials are available, use them
	if cfg.MongoDB.Username != "" {
		// Construct connection string
		connectionString = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?%s",
			cfg.MongoDB.Username,
			cfg.MongoDB.Password,
			cfg.MongoDB.Host,
			cfg.MongoDB.Database,
			cfg.MongoDB.ConnectionOptions)
	} else {
		// Provide a fallback for testing
		connectionString = fmt.Sprintf("mongodb://%s", cfg.MongoDB.Host)
		log.Println("Using connection string: " + connectionString)
	}

	// Create client and connect
	clientOptions := options.Client().ApplyURI(connectionString)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDB.Timeout)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Verify connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}
