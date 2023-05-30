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
	// reddit config
	ClientId     string `env:"CLIENT_ID,notEmpty"`
	ClientSecret string `env:"CLIENT_SECRET,notEmpty"`
	Username     string `env:"BOT_USERNAME,notEmpty"`
	Password     string `env:"BOT_PASSWORD,notEmpty"`
	Subreddit    string `env:"SUBREDDIT,notEmpty"`

	// techpowerup config
	TPUHost     string `env:"TPU_HOST,notEmpty"`
	TPUUsername string `env:"TPU_USERNAME,notEmpty"`
	TPUSecret   string `env:"TPU_SECRET,notEmpty"`

	// elasticsearch config
	EsAddress string `env:"ES_ADDRESS,notEmpty"`

	// application config
	OverrideOldBot bool `env:"OVERRIDE_OLD_BOT,notEmpty"`

	//debugging config
	Token           string `env:"BOT_ACCESS_TOKEN"`
	ExpireTimeMilli int64  `env:"BOT_TOKEN_EXPIRE_MILLI"`
	IsDebug         bool   `env:"IS_DEBUG"`
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
