package main

import (
	"context"
	"time"

	"github.com/aattwwss/ssd-bot-go/elasticutil"
	"github.com/aattwwss/ssd-bot-go/pkg/reddit"
	"github.com/aattwwss/ssd-bot-go/pkg/ssd"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	config := config{}
	if err := env.Parse(&config); err != nil {
		log.Fatal().Msgf("Parse env error: %v", err)
	}

	_, err = reddit.NewRedditClient(config.ClientId, config.ClientSecret, config.Username, config.Password, config.Token, config.ExpireTimeMilli, config.OverrideOldBot, config.IsDebug)
	if err != nil {
		log.Fatal().Msgf("Init reddit client error: %v", err)
	}
	es, _ := elasticutil.NewElasticsearchClient(config.EsAddress)
	esRepo := ssd.NewEsSSDRepository(es, "ssd-index")
	tpuRepo := ssd.NewTpuSSDRepository(config.TPUHost, config.TPUUsername, config.TPUSecret)
	ssdSync := ssd.SSDSync{
		StartId: 40,
		// EndId:    50,
		Delay:    time.Duration(100),
		IdToSkip: []int{},
	}
	err = ssdSync.Sync(context.Background(), tpuRepo, esRepo)
	if err != nil {
		log.Fatal().Msgf("sync error", err)
	}
	// ssd, _ := esRepo.FindById(context.Background(), "123")
	// ssds, _ := esRepo.SearchBasic(context.Background(), "corsair")
	// sss, _ := esRepo.Search(context.Background(), "corsair")
	// log.Info().Msgf("%v", ssd)
	// log.Info().Msgf("%v", ssds)
	// log.Info().Msgf("%v", sss)
	// ssd.DriveID = ssd.DriveID + "_new"
	// ssd.Capacity = "some capacity"
	// esRepo.Insert(context.Background(), *ssd)
}

type config struct {
	// reddit config
	ClientId     string `env:"CLIENT_ID,notEmpty"`
	ClientSecret string `env:"CLIENT_SECRET,notEmpty"`
	Username     string `env:"BOT_USERNAME,notEmpty"`
	Password     string `env:"BOT_PASSWORD,notEmpty"`

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
