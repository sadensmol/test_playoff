package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sadensmol/test_playoff/internal/config"
	"github.com/sadensmol/test_playoff/internal/invite"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type App struct {
	cfg config.Config
}

func NewApp(cfg config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) initDB(ctx context.Context, cfg config.Mongo) (*mongo.Client, error) {
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

func (a *App) Run(ctx context.Context) {

	cl, err := a.initDB(ctx, a.cfg.Mongo)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}

	invService := invite.NewInviteService(cl)

	e := echo.New()
	invite.NewHandler(invService).Register(e)

	intCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", a.cfg.HTTP.Port)); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	<-intCtx.Done()
	intCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := e.Shutdown(intCtx); err != nil {
		log.Error().Err(err).Msg("error while shutting down the server")
	}

}
