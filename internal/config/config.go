package config

import (
	"fmt"
	"sync"

	"github.com/codingconcepts/env"
	"github.com/rs/zerolog/log"
)

type Config struct {
	HTTP
	Mongo
}

type HTTP struct {
	Port int `env:"HTTP_PORT"`
}

type Mongo struct {
	Host          string `env:"MONGODB_HOST"`
	Port          int    `env:"MONGODB_PORT"`
	User          string `env:"MONGODB_USER"`
	AdminDatabase string `env:"MONGODB_ADMIN_DATABASE"`
	Password      string `env:"MONGODB_PASSWORD"`
	MaxPoolSize   int    `env:"MONGODB_MAX_POOL_SIZE"`
}

func (m Mongo) ConnectionString() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?maxPoolSize=%d&connectTimeoutMS=5000&socketTimeoutMS=5000",
		m.User,
		m.Password,
		m.Host,
		m.Port,
		m.AdminDatabase,
		m.MaxPoolSize)
}

var once sync.Once
var instance Config

func GetConfig() Config {
	once.Do(func() {
		if err := env.Set(&instance); err != nil {
			log.Fatal().Msgf("cannot init config %s", err)
		}
		if err := env.Set(&instance.HTTP); err != nil {
			log.Fatal().Msgf("cannot init config %s", err)
		}
		if err := env.Set(&instance.Mongo); err != nil {
			log.Fatal().Msgf("cannot init config %s", err)
		}
	})

	return instance
}
