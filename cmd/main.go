package main

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/sadensmol/test_playoff/internal"
	"github.com/sadensmol/test_playoff/internal/config"
	"github.com/sadensmol/test_playoff/internal/utils"
)

func main() {
	utils.InitLogger(zerolog.InfoLevel)

	internal.NewApp(config.GetConfig()).Run(context.Background())
}
