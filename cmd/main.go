package main

import (
	"github.com/aattwwss/ssd-bot-go/pkg/reddit"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	SUBREDDIT   = "buildapcsales"
	LINK_PREFIX = "t3_"

	SPREADSHEET_ID = "1B27_j9NDPU3cNlj2HKcrfpJKHkOf-Oi1DbuuQva2gT4"
	SHEET_NAME     = "'Master List'" //take note of the single quote, which is needed for sheets with space in them
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

	rc, err := reddit.NewRedditClient(config.ClientId, config.ClientSecret, config.Username, config.Password, config.Token, config.ExpireTimeMilli, config.IsDebug)
	if err != nil {
		log.Error().Msgf("Init reddit client error: %v", err)
		return
	}

}

type config struct {
	ClientId       string `env:"CLIENT_ID,notEmpty"`
	ClientSecret   string `env:"CLIENT_SECRET,notEmpty"`
	Username       string `env:"BOT_USERNAME,notEmpty"`
	Password       string `env:"BOT_PASSWORD,notEmpty"`
	TPUHost        string `env:"TPU_HOST,notEmpty"`
	TPUSecret      string `env:"TPU_SECRET"`
	EsAccessKey    string `env:"ES_ACCESS_KEY"`
	EsAccessSecret string `env:"ES_ACCESS_SECRET"`

	Token           string `env:"BOT_ACCESS_TOKEN"`
	ExpireTimeMilli int64  `env:"BOT_TOKEN_EXPIRE_MILLI"`
	IsDebug         bool   `env:"IS_DEBUG"`
}
