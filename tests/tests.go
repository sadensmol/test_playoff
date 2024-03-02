package tests

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sadensmol/test_playoff/internal"
	"github.com/sadensmol/test_playoff/internal/config"
	"github.com/sadensmol/test_playoff/internal/utils"
)

type IntegrationTest struct {
	Cfg        config.Config
	wg         *sync.WaitGroup
	cancelFunc context.CancelFunc
}

func (s *IntegrationTest) TearDown() {
	s.cancelFunc()
	s.wg.Wait()
	time.Sleep(1 * time.Second) // wait for app to stop
}

func (s *IntegrationTest) Setup() {
	httpPort := utils.GetRandomUnusedPort()

	s.Cfg = config.Config{
		HTTP: config.HTTP{Port: httpPort},
		Mongo: config.Mongo{
			Host:          "localhost",
			Port:          27017,
			User:          "mongo",
			Password:      "mongo",
			AdminDatabase: "admin",
			MaxPoolSize:   100,
		},
	}

	s.wg = &sync.WaitGroup{}
	s.wg.Add(1)

	go func() {
		utils.InitLogger(zerolog.DebugLevel)
		output := zerolog.ConsoleWriter{Out: os.Stderr}
		log.Logger = log.Output(output)

		log.Info().Msg("starting app")
		ctx, cancel := context.WithCancel(context.Background())
		s.cancelFunc = cancel
		internal.NewApp(s.Cfg).Run(ctx)
		log.Info().Msg("app stopped!")
		s.wg.Done()
	}()

	time.Sleep(1 * time.Second) // wait for app to start
}
