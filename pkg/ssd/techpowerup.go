package ssd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type tpuResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  T      `json:"result"`
}

type TpuSSDRepository struct {
	host     string
	username string
	apikey   string
}

func NewTpuSSDRepository(host, username, apiKey string) *TpuSSDRepository {
	return &TpuSSDRepository{
		host:     host,
		username: username,
		apikey:   apiKey,
	}
}

func (tpu *TpuSSDRepository) buildUrl() string {
	return fmt.Sprintf("%s/ssd-specs/api/%s/v1", tpu.host, tpu.username)
}

func (tpu *TpuSSDRepository) FindById(ctx context.Context, id string) (*SSD, error) {
	url := tpu.buildUrl() + fmt.Sprintf("/query?key=%s&id=%s", tpu.apikey, id)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var tpuRes tpuResponse[SSD]
	err = json.NewDecoder(response.Body).Decode(&tpuRes)
	if err != nil {
		return nil, err
	}
	if tpuRes.Status == "failed" && tpuRes.Message == "Drive not found" {
		return nil, nil
	}
	if tpuRes.Status != "success" {
		fmt.Printf("%v", response.StatusCode)
		fmt.Printf("%v", tpuRes.Status)
		return nil, errors.New("tpu query status error: " + tpuRes.Status)
	}
	return &tpuRes.Result, nil
}

func (tpu *TpuSSDRepository) SearchBasic(ctx context.Context, s string) ([]BasicSSD, error) {
	url := tpu.buildUrl() + fmt.Sprintf("/lookup?key=%s&id=%s", tpu.apikey, s)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var tpuRes tpuResponse[[]BasicSSD]
	err = json.NewDecoder(response.Body).Decode(&tpuRes)
	if err != nil {
		return nil, err
	}
	if tpuRes.Status == "failed" && tpuRes.Message == "Drive not found" {
		return nil, nil
	}
	if tpuRes.Status != "success" {
		return nil, errors.New("tpu lookup status error: " + err.Error())
	}
	return tpuRes.Result, nil
}

func (tpu *TpuSSDRepository) Search(ctx context.Context, s string) ([]SSD, error) {
	return nil, nil
}

func (tpu *TpuSSDRepository) Insert(ctx context.Context, ssd SSD) error {
	//TODO implement this
	return nil
}

func (tpu *TpuSSDRepository) Update(ctx context.Context, ssd SSD) error {
	//TODO implement this
	return nil
}
