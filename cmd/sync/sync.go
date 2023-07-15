package main

import (
	"context"
	"flag"
	"github.com/aattwwss/ssd-bot-go/internal/config"
	"github.com/aattwwss/ssd-bot-go/pkg/ssd"
	"strconv"
	"time"

	"github.com/aattwwss/ssd-bot-go/elasticutil"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type syncParam struct {
	StartId  int
	EndId    int
	IdToSkip []int
	Delay    time.Duration
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatal().Msgf("Parse env error: %v", err)
	}
	// Define command-line flags
	startId := flag.Int("startId", 1, "Start ID to sync from")
	endId := flag.Int("endId", 1550, "End ID to sync to")
	flag.Parse()

	es, _ := elasticutil.NewElasticsearchClient(cfg.EsAddress)
	esRepo := ssd.NewEsRepository(es, "ssd-index")
	tpuRepo := ssd.NewTpuRepository(cfg.TPUHost, cfg.TPUUsername, cfg.TPUSecret)

	param := syncParam{
		StartId:  *startId,
		EndId:    *endId,
		Delay:    time.Duration(10),
		IdToSkip: []int{},
	}
	err = sync(context.Background(), tpuRepo, esRepo, param)
	if err != nil {
		log.Fatal().Msgf("sync error", err)
	}
}

func sync(ctx context.Context, source ssd.Repository, destination ssd.Repository, s syncParam) error {
	for id := s.StartId; id <= s.EndId; id++ {
		if contains(s.IdToSkip, id) {
			log.Info().Msgf("Skipping id: %v", id)
			continue
		}
		log.Info().Msgf("Syncing with id: %v", id)
		found, err := source.FindById(ctx, strconv.Itoa(id))

		if err != nil {
			log.Error().Msgf("Source find by id, id: %v, error: %v", id, err)
			return err
		}
		if found == nil {
			log.Info().Msgf("Source find by id returns empty with id: %v", id)
			continue
			// return nil
		}
		err = destination.Insert(ctx, *found)
		if err != nil {
			log.Error().Msgf("Destination insert by id, id: %v, error: %v", id, err)
			continue
		}
		time.Sleep(s.Delay)
	}
	return nil
}

func contains[T comparable](arr []T, element T) bool {
	for _, item := range arr {
		if item == element {
			return true
		}
	}
	return false
}
