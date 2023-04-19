package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/aattwwss/ssd-bot-go/benchmark"
	"github.com/aattwwss/ssd-bot-go/bot"
	"github.com/aattwwss/ssd-bot-go/config"
	"github.com/aattwwss/ssd-bot-go/reddit"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	config := config.Config{}
	if err := env.Parse(&config); err != nil {
		log.Fatal().Msgf("Parse env error: %v", err)
	}

	// Define a flag called "benchmarkFlag" that accepts a boolean value and has a default value of false
	benchmarkFlag := flag.Bool("benchmark", false, "run benchmark instead of running the bot indefinitely")

	// Parse the command line arguments
	flag.Parse()
	if *benchmarkFlag {
		fmt.Println("Running benchmark...")
		benchmark.Benchmark(config)
	} else {
		rc, err := reddit.NewRedditClient(config.ClientId, config.ClientSecret, config.Username, config.Password, config.Token, config.ExpireTimeMilli, config.IsDebug)
		if err != nil {
			log.Error().Msgf("Init reddit client error: %v", err)
			return
		}

		for {
			log.Info().Msgf("Scanning...")
			count, err := bot.Run(config, rc)
			if err != nil {
				log.Error().Msgf("Run error: %v", err)
			}
			log.Info().Msgf("Updated %v posts...", count)
			time.Sleep(10 * time.Minute)
		}
	}
}
