package ssd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/aattwwss/ssd-bot-go/elasticutil"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/rs/zerolog/log"
)

type EsRepository struct {
	EsClient *elasticsearch.Client
	Index    string
}

func NewEsRepository(esClient *elasticsearch.Client, index string) *EsRepository {
	return &EsRepository{
		EsClient: esClient,
		Index:    index,
	}
}

func (esRepo *EsRepository) FindById(ctx context.Context, driveId string) (*SSD, error) {
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

func (esRepo *EsRepository) SearchBasic(ctx context.Context, s string) ([]SSDBasic, error) {
	var ssdResponse elasticutil.SearchResponse[SSDBasic]
	var res []SSDBasic
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

// BoolQuery Elastic bool query
type BoolQuery struct {
	Bool BoolQueryParams `json:"bool"`
}

// BoolQueryParams params for an Elastic bool query
type BoolQueryParams struct {
	Must               []interface{} `json:"must,omitempty"`
	Should             []interface{} `json:"should,omitempty"`
	Filter             []interface{} `json:"filter,omitempty"`
	MinimumShouldMatch int           `json:"minimum_should_match,omitempty"`
}

func (esRepo *EsRepository) Search(ctx context.Context, searchQuery string) ([]SSD, error) {
	var ssdResponse elasticutil.SearchResponse[SSD]
	var res []SSD

	boolQuery := BoolQuery{
		Bool: BoolQueryParams{
			Must: []interface{}{},
		},
	}

	matchQuery := map[string]interface{}{
		"multi_match": map[string]interface{}{
			"query": searchQuery,
		},
	}
	boolQuery.Bool.Must = append(boolQuery.Bool.Must, matchQuery)

	capacity, ok := parseCapacity(searchQuery)
	if ok {
		capacityQuery := map[string]interface{}{
			"term": map[string]interface{}{
				"capacity": capacity,
			},
		}
		boolQuery.Bool.Must = append(boolQuery.Bool.Must, capacityQuery)
	}

	ff, ok := parseFormFactor(searchQuery)
	if ok {
		ffQuery := map[string]interface{}{
			"term": map[string]interface{}{
				"formFactor": ff,
			},
		}
		boolQuery.Bool.Must = append(boolQuery.Bool.Must, ffQuery)
	}

	query := map[string]interface{}{
		"query": boolQuery,
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

func (esRepo *EsRepository) Update(ctx context.Context, ssd SSD) error {
	//TODO implement this
	return nil
}

func (esRepo *EsRepository) Insert(ctx context.Context, ssd SSD) error {
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
		log.Error().Msgf("[%s] Error indexing document ID=%v", res.Status(), ssd.DriveID)
		return errors.New("index ssd response error")
	}
	return nil
}

// wrapper for search queries to Elastic client
func (esRepo *EsRepository) doSearch(ctx context.Context, query map[string]interface{}, payload any) error {
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

// parseCapacity parses a string for a capacity in TB or GB
func parseCapacity(s string) (int, bool) {
	s = strings.ToUpper(s)
	re := regexp.MustCompile(`(\d+)\s*(TB|GB)`)
	match := re.FindStringSubmatch(s)
	if len(match) <= 1 {
		return 0, false
	}
	capacity, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, false
	}
	return capacity, true
}

// parseFormFactor parses a string for the length of a ssd
func parseFormFactor(s string) (int, bool) {
	ssdLength := []int{
		2230,
		2280,
		2242,
		22110,
	}
	for _, l := range ssdLength {
		if strings.Contains(s, strconv.Itoa(l)) {
			return l, true
		}
	}
	return 0, false
}
