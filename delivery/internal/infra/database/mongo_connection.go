package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewMongoConnection(dns string) (*mongo.Client, error) {
	const maxAttempts = 5

	var lastErr error

	for i := 0; i < maxAttempts; i++ {
		clientOptions := options.Client().ApplyURI(dns)

		client, err := mongo.Connect(clientOptions)
		if err != nil {
			lastErr = err
		} else {
			ctx, cancel := context.WithTimeout(
				context.Background(),
				5*time.Second,
			)

			err = client.Ping(ctx, nil)
			cancel()

			if err == nil {
				return client, nil
			}

			lastErr = err
			_ = client.Disconnect(context.Background())
		}

		if i < maxAttempts-1 {
			time.Sleep(2 * time.Second)
		}
	}

	return nil, fmt.Errorf(
		"unable to connect to MongoDB after %d attempts: %w",
		maxAttempts,
		lastErr,
	)
}
