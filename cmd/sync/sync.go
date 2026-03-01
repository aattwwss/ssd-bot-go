package main

import (
	"context"
	"flag"
	"slices"

	"github.com/aattwwss/ssd-bot-go/internal/config"
	"github.com/aattwwss/ssd-bot-go/pkg/ssd"
	"strconv"
	"time"

	"github.com/aattwwss/ssd-bot-go/elasticutil"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const (
	ES_INDEX          = "ssd-index"
	DEFAULT_START_ID  = 1
	DEFAULT_END_ID    = 1550
	DEFAULT_SYNC_DELAY = 10 * time.Second
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
	startId := flag.Int("startId", DEFAULT_START_ID, "Start ID to sync from")
	endId := flag.Int("endId", DEFAULT_END_ID, "End ID to sync to")
	flag.Parse()

	es, err := elasticutil.NewElasticsearchClient(cfg.EsAddress)
	if err != nil {
		log.Fatal().Msgf("Init elasticsearch client error: %v", err)
	}
	esRepo := ssd.NewEsRepository(es, ES_INDEX)
	tpuRepo := ssd.NewTpuRepository(cfg.TPUHost, cfg.TPUUsername, cfg.TPUSecret)

	param := syncParam{
		StartId:  *startId,
		EndId:    *endId,
		Delay:    DEFAULT_SYNC_DELAY,
		IdToSkip: nil,
	}
	err = sync(context.Background(), tpuRepo, esRepo, param)
	if err != nil {
		log.Fatal().Err(err).Msg("Sync error")
	}
}

func sync(ctx context.Context, source ssd.Repository, destination ssd.Repository, s syncParam) error {
	for id := s.StartId; id <= s.EndId; id++ {
		if slices.Contains(s.IdToSkip, id) {
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
