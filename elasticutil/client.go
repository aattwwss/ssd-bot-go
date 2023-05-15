package elasticutil

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

type SearchResponse[T any] struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string  `json:"_index"`
			Type   string  `json:"_type"`
			ID     string  `json:"_id"`
			Score  float64 `json:"_score"`
			Source T       `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
