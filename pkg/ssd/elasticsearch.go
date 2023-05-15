package ssd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/aattwwss/ssd-bot-go/elasticutil"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/rs/zerolog/log"
)

type EsSSDRepository struct {
	EsClient *elasticsearch.Client
	Index    string
}

func NewEsSSDRepository(esClient *elasticsearch.Client, index string) *EsSSDRepository {
	return &EsSSDRepository{
		EsClient: esClient,
		Index:    index,
	}
}

func (esRepo *EsSSDRepository) FindById(ctx context.Context, driveId string) (*SSD, error) {
	var ssdResponse elasticutil.SearchResponse[SSD]
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"driveId": driveId,
			},
		},
	}
	err := esRepo.doSearch(ctx, query, &ssdResponse)
	if err != nil {
		return nil, err
	}
	if len(ssdResponse.Hits.Hits) == 0 {
		return nil, nil
	}
	return &ssdResponse.Hits.Hits[0].Source, nil
}

func (esRepo *EsSSDRepository) Search(context context.Context, s string) ([]BasicSSD, error) {
	//TODO implement this
	return nil, nil
}

func (esRepo *EsSSDRepository) Insert(context context.Context, ssd SSD) error {
	//TODO implement this
	return nil
}

func (esRepo *EsSSDRepository) Update(context context.Context, ssd SSD) error {
	//TODO implement this
	return nil
}

func (esRepo *EsSSDRepository) doSearch(ctx context.Context, query map[string]interface{}, payload any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Error().Msgf("Error encoding query: %s", err)
		return errors.New("error decoding query")
	}
	log.Info().Msgf("asa" + buf.String())
	es := esRepo.EsClient
	res, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex(esRepo.Index),
		es.Search.WithBody(&buf),
	)
	if err != nil {
		log.Error().Msgf("Error getting response: %s", err)
		return errors.New("search response error")
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Error().Msgf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Error().Msgf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
		return errors.New("find by id response error")
	}

	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		log.Error().Msgf("Error parsing the response body: %s", err)
		return errors.New("search payload decode error")
	}
	return nil
}
