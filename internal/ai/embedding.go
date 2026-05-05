package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type OpenAIEmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

type OpenAIEmbeddingResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}

func GetEmbedding(ctx context.Context, apiKey string, text string) ([]float32, error) {
	apiURL := "https://api.openai.com/v1/embeddings"

	transport := &http.Transport{
		Proxy: nil,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	reqBody, _ := json.Marshal(OpenAIEmbeddingRequest{
		Input: text,
		Model: "text-embedding-3-small",
	})

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI-ga ulanishda xato: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI xatosi, status kodi: %d", resp.StatusCode)
	}

	var result OpenAIEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data[0].Embedding, nil
}
