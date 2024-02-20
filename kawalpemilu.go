package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type KawalPemiluClient struct {
	BaseURL string
}

func (c *KawalPemiluClient) GetNationalCountingResult(ctx context.Context) (KawalPemiluResponse, error) {
	return c.GetCountingResultByID(ctx, "")
}

func (c *KawalPemiluClient) GetAllProvincesCountingResult(ctx context.Context) (KawalPemiluResponse, error) {
	result := KawalPemiluResponse{
		Result: KawalPemiluResponseData{
			Aggregated: make(map[string][]CountingResult, len(provinceByID)),
		},
	}

	for provinceID, province := range provinceByID {
		provinceRes, err := c.GetCountingResultByID(ctx, provinceID)
		if err != nil {
			return KawalPemiluResponse{}, fmt.Errorf("get province %s counting result: %w", province, err)
		}

		for id, res := range provinceRes.Result.Aggregated {
			result.Result.Aggregated[id] = append(result.Result.Aggregated[id], res...)
		}
	}

	return result, nil
}

func (c *KawalPemiluClient) GetCountingResultByID(ctx context.Context, id string) (KawalPemiluResponse, error) {
	path := fmt.Sprintf("%s/h?id=%s", c.BaseURL, id)
	req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
	if err != nil {
		return KawalPemiluResponse{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	log.Printf("Calling GET %s", path)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return KawalPemiluResponse{}, fmt.Errorf("do request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return KawalPemiluResponse{}, fmt.Errorf("got status %d [%s]", res.StatusCode, res.Status)
	}

	log.Printf("Success GET %s", path)

	resBytes, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return KawalPemiluResponse{}, fmt.Errorf("read body: %w", err)
	}

	var response KawalPemiluResponse
	if err := json.Unmarshal(resBytes, &response); err != nil {
		return KawalPemiluResponse{}, fmt.Errorf("unmarshal response: %w", err)
	}

	return response, nil
}
