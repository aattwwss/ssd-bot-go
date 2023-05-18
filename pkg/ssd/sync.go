package ssd

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
)

type SSDSynchronizer interface {
	Sync(ctx context.Context, source SSDRepository, destination SSDRepository) error
}

type SSDSync struct {
	StartId  int
	EndId    int
	IdToSkip []int
	Delay    time.Duration
}

func (s SSDSync) Sync(ctx context.Context, source SSDRepository, destination SSDRepository) error {
	if s.EndId == 0 {
		s.EndId = 1490
	}

	for id := s.StartId; id <= s.EndId; id++ {
		if contains(s.IdToSkip, id) {
			log.Info().Msgf("Skipping id: %v", id)
			continue
		}
		log.Info().Msgf("Syncing with id: %v", id)
		ssd, err := source.FindById(ctx, strconv.Itoa(id))

		if err != nil {
			log.Error().Msgf("Source find by id, id: %v, error: %v", id, err)
			return err
		}
		if ssd == nil {
			log.Info().Msgf("Source find by id returns empty with id: %v", id)
			continue
			// return nil
		}
		err = destination.Insert(ctx, *ssd)
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
