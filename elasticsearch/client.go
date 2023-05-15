package elasticsearch

import (
	"errors"
	"fmt"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/rs/zerolog/log"
)

func NewElasticsearchClient(address string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{address},
	}
	es, err := elasticsearch.NewClient(cfg)
	res, err := es.Ping()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(fmt.Sprintf("es ping error: %s", res.String()))
	} else {
		log.Info().Msg("Elasticsearch connection successful!")
	}
	return es, nil
}
