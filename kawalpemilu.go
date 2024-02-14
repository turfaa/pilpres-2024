package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type KawalPemiluClient struct {
	BaseURL string
}

func (c *KawalPemiluClient) GetNationalCountingResult(ctx context.Context) (KawalPemiluResponse, error) {
	js, err := json.Marshal(KawalPemiluRequest{})
	if err != nil {
		return KawalPemiluResponse{}, fmt.Errorf("marshal json: %w", err)
	}

	jsBuf := bytes.NewBuffer(js)

	path := fmt.Sprintf("%s/hierarchy2", c.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", path, jsBuf)
	if err != nil {
		return KawalPemiluResponse{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return KawalPemiluResponse{}, fmt.Errorf("do request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return KawalPemiluResponse{}, fmt.Errorf("got status %d [%s]", res.StatusCode, res.Status)
	}

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
