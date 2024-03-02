package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sadensmol/test_playoff/internal/config"
	"github.com/sadensmol/test_playoff/internal/invite"
)

type App struct {
	cfg config.Config
}

func NewApp(cfg config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

func (a *App) Run(ctx context.Context) {
	cl, err := InitMongo(ctx, a.cfg.Mongo)
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
