package ssd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/aattwwss/ssd-bot-go/elasticutil"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
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

func (esRepo *EsSSDRepository) SearchBasic(ctx context.Context, s string) ([]BasicSSD, error) {
	//TODO implement this
	var ssdResponse elasticutil.SearchResponse[BasicSSD]
	var res []BasicSSD
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": s,
			},
		},
	}
	err := esRepo.doSearch(ctx, query, &ssdResponse)
	if err != nil {
		return nil, err
	}

	for _, hit := range ssdResponse.Hits.Hits {
		res = append(res, hit.Source)
	}
	return res, nil
}

func (esRepo *EsSSDRepository) Search(ctx context.Context, searchQuery string) ([]SSD, error) {
	//TODO implement this
	var ssdResponse elasticutil.SearchResponse[SSD]
	var res []SSD
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query": searchQuery,
			},
		},
	}
	err := esRepo.doSearch(ctx, query, &ssdResponse)
	if err != nil {
		return nil, err
	}

	for _, hit := range ssdResponse.Hits.Hits {
		if sanityCheck(searchQuery, hit.Source) {
			res = append(res, hit.Source)
		}
	}
	return res, nil
}

// a sanityCheck to ensure we only return the ssd we know is correct
// remove false positive as much as possible
func sanityCheck(searchQuery string, ssd SSD) bool {
	if !strings.Contains(strings.ToLower(strings.ReplaceAll(searchQuery, " ", "")), strings.ToLower(strings.ReplaceAll(ssd.Manufacturer, " ", ""))) {
		return false
	}
	cleanedName := cleanName(ssd.Name)
	if !strings.Contains(strings.ToLower(strings.ReplaceAll(searchQuery, " ", "")), strings.ToLower(strings.ReplaceAll(cleanedName, " ", ""))) {
		return false
	}
	return true
}

// cleanName will remove any redundant characters so that the sanity check can
// be more accurate. Ideally we would want to do this in our persistence layer,
// but I'm lazy
func cleanName(name string) string {
	return strings.ReplaceAll(name, "(w/ Heatsink)", "")
}

func (esRepo *EsSSDRepository) Update(ctx context.Context, ssd SSD) error {
	//TODO implement this
	return nil
}

func (esRepo *EsSSDRepository) Insert(ctx context.Context, ssd SSD) error {
	// Build the request body.
	data, err := json.Marshal(ssd)
	if err != nil {
		log.Error().Msgf("Error marshaling document: %s", err)
		return err
	}

	// Set up the request object.
	req := esapi.IndexRequest{
		Index:      esRepo.Index,
		DocumentID: ssd.DriveID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(ctx, esRepo.EsClient)
	if err != nil {
		log.Error().Msgf("Error getting response: %s", err)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Error().Msgf("[%s] Error indexing document ID=%d", res.Status(), ssd.DriveID)
		return errors.New("index ssd response error")
	}
	return nil
}

func (esRepo *EsSSDRepository) doSearch(ctx context.Context, query map[string]interface{}, payload any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Error().Msgf("Error encoding query: %s", err)
		return errors.New("error decoding query")
	}
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
