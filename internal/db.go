package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sadensmol/test_playoff/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongo(ctx context.Context, cfg config.Mongo) (*mongo.Client, error) {
	mongoCtx, cancelFunc := context.WithTimeout(ctx, time.Second*5)
	defer cancelFunc()

	client, err := mongo.Connect(mongoCtx, options.Client().ApplyURI(cfg.ConnectionString()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().Msg("database initialized")

	return client, nil
}
