package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type KawalPemiluClient struct {
	BaseURL     string
	Concurrency int
}

func (c *KawalPemiluClient) GetNationalCountingResult(ctx context.Context) (KawalPemiluResponse, error) {
	return c.GetCountingResultByID(ctx, "")
}

func (c *KawalPemiluClient) GetAllProvincesCountingResult(ctx context.Context) (KawalPemiluResponse, error) {
	ids := make([]string, 0, len(provinceByID))
	for id := range provinceByID {
		ids = append(ids, id)
	}

	return c.GetCountingResultByIDs(ctx, ids)
}

func (c *KawalPemiluClient) GetAllCitiesCountingResult(ctx context.Context) (KawalPemiluResponse, error) {
	ids := make([]string, 0, len(cityByID))
	for id := range cityByID {
		ids = append(ids, id)
	}

	return c.GetCountingResultByIDs(ctx, ids)
}

func (c *KawalPemiluClient) GetCountingResultByIDs(ctx context.Context, ids []string) (KawalPemiluResponse, error) {
	result := KawalPemiluResponse{
		Result: KawalPemiluResponseData{
			Aggregated: make(map[string][]CountingResult, len(ids)),
		},
	}

	var (
		wg          sync.WaitGroup
		lock        sync.Mutex
		concurrency = make(chan struct{}, c.Concurrency)
	)

	for _, id := range ids {
		wg.Add(1)

		go func(id string) {
			defer wg.Done()

			concurrency <- struct{}{}
			defer func() {
				<-concurrency
			}()

			res, err := c.GetCountingResultByID(ctx, id)
			if err != nil {
				log.Printf("error get place %s counting result: %s", id, err)
				return
			}

			lock.Lock()
			for resID, countingResults := range res.Result.Aggregated {
				result.Result.Aggregated[resID] = append(result.Result.Aggregated[resID], countingResults...)
			}
			lock.Unlock()
		}(id)
	}

	wg.Wait()
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
