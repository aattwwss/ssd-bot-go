package main

import (
	"context"
	"time"

	"github.com/aattwwss/ssd-bot-go/elasticutil"
	"github.com/aattwwss/ssd-bot-go/pkg/ssd"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type config struct {
	// techpowerup config
	TPUHost     string `env:"TPU_HOST,notEmpty"`
	TPUUsername string `env:"TPU_USERNAME,notEmpty"`
	TPUSecret   string `env:"TPU_SECRET,notEmpty"`

	// elasticsearch config
	EsAddress string `env:"ES_ADDRESS,notEmpty"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	config := config{}
	if err := env.Parse(&config); err != nil {
		log.Fatal().Msgf("Parse env error: %v", err)
	}

	es, _ := elasticutil.NewElasticsearchClient(config.EsAddress)
	esRepo := ssd.NewEsSSDRepository(es, "ssd-index")
	tpuRepo := ssd.NewTpuSSDRepository(config.TPUHost, config.TPUUsername, config.TPUSecret)
	sync(tpuRepo, esRepo)
}

func sync(source, dest ssd.SSDRepository) {
	ssdSync := ssd.SSDSync{
		StartId:  1,
		EndId:    1520,
		Delay:    time.Duration(10),
		IdToSkip: []int{},
	}
	err := ssdSync.Sync(context.Background(), source, dest)
	if err != nil {
		log.Fatal().Msgf("sync error", err)
	}
}
