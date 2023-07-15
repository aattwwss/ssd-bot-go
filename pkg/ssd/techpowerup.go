package ssd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type response[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  T      `json:"result"`
}

type TpuRepository struct {
	host     string
	username string
	apikey   string
}

func NewTpuRepository(host, username, apiKey string) *TpuRepository {
	return &TpuRepository{
		host:     host,
		username: username,
		apikey:   apiKey,
	}
}

func (tpu *TpuRepository) FindById(ctx context.Context, id string) (*SSD, error) {
	url := fmt.Sprintf("%s/ssd-specs/api/%s/v1/query?key=%s&id=%s", tpu.host, tpu.username, tpu.apikey, id)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tpuRes response[SSD]
	err = json.NewDecoder(resp.Body).Decode(&tpuRes)
	if err != nil {
		return nil, err
	}
	if tpuRes.Status == "failed" && tpuRes.Message == "Drive not found" {
		return nil, nil
	}
	if tpuRes.Status != "success" {
		fmt.Printf("%v", resp.StatusCode)
		fmt.Printf("%v", tpuRes.Status)
		return nil, errors.New("tpu query status error: " + tpuRes.Status)
	}
	return &tpuRes.Result, nil
}

func (tpu *TpuRepository) SearchBasic(ctx context.Context, s string) ([]SSDBasic, error) {
	url := fmt.Sprintf("%s/ssd-specs/api/%s/v1/lookup?key=%s&id=%s", tpu.host, tpu.username, tpu.apikey, s)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tpuRes response[[]SSDBasic]
	err = json.NewDecoder(resp.Body).Decode(&tpuRes)
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

func (tpu *TpuRepository) Search(ctx context.Context, s string) ([]SSD, error) {
	return nil, nil
}

func (tpu *TpuRepository) Insert(ctx context.Context, ssd SSD) error {
	//TODO implement this
	return nil
}

func (tpu *TpuRepository) Update(ctx context.Context, ssd SSD) error {
	//TODO implement this
	return nil
}
