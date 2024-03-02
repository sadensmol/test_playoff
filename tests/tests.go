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
	"go.mongodb.org/mongo-driver/mongo"
)

type IntegrationTest struct {
	MongoClient *mongo.Client
	Ctx         *context.Context
	Cfg         config.Config
	wg          *sync.WaitGroup
	cancelFunc  context.CancelFunc
}

func (s *IntegrationTest) TearDown() {
	s.cancelFunc()
	s.wg.Wait()
	time.Sleep(1 * time.Second) // wait for app to stop
}

func (s *IntegrationTest) Setup() {
	httpPort := utils.GetRandomUnusedPort()
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel
	s.Ctx = &ctx

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

	cl, err := internal.InitMongo(ctx, s.Cfg.Mongo)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize database")
	}

	s.MongoClient = cl

	s.wg = &sync.WaitGroup{}
	s.wg.Add(1)

	go func() {
		utils.InitLogger(zerolog.DebugLevel)
		output := zerolog.ConsoleWriter{Out: os.Stderr}
		log.Logger = log.Output(output)
		log.Info().Msg("starting app")
		internal.NewApp(s.Cfg).Run(ctx)
		log.Info().Msg("app stopped!")
		s.wg.Done()
	}()

	time.Sleep(1 * time.Second) // wait for app to start
}
