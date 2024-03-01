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
)

type App struct {
	cfg config.Config
}

func NewApp(cfg config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) setupDB(cfg config.Mongo) {
	log.Info().Msg("database initialized")
}

func (a *App) Run(ctx context.Context) {

	a.setupDB(a.cfg.Mongo)

	e := echo.New()
	invite.NewHandler().Register(e)

	intCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	go func() {
		if err := e.Start(fmt.Sprintf(":%d", a.cfg.HTTP.Port)); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("shutting down the server")
		}
	}()

	<-intCtx.Done()
	intCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := e.Shutdown(intCtx); err != nil {
		log.Error().Err(err).Msg("error while shutting down the server")
	}

}
