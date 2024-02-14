package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type Service struct {
	KawalPemiluClient *KawalPemiluClient
	Predictor         *SimplePredictor

	currentPrediction CountingResult
	lock              sync.RWMutex
}

func (s *Service) RunRefresher(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		if err := s.RefreshPrediction(ctx); err != nil {
			log.Printf("error refreshing prediction: %s", err)
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) RefreshPrediction(ctx context.Context) error {
	res, err := s.KawalPemiluClient.GetNationalCountingResult(ctx)
	if err != nil {
		return fmt.Errorf("get national counting result from client: %w", err)
	}

	prediction := s.Predictor.PredictResult(extractResults(res.Result))

	s.lock.Lock()
	s.currentPrediction = prediction
	s.lock.Unlock()

	return nil
}

func extractResults(data KawalPemiluResponseData) []CountingResult {
	var results []CountingResult
	for _, r := range data.Aggregated {
		results = append(results, r...)
	}
	return results
}

func (s *Service) GetNationalCountingPrediction(_ context.Context) (CountingResult, error) {
	s.lock.RLock()
	result := s.currentPrediction
	s.lock.RUnlock()
	return result, nil
}
