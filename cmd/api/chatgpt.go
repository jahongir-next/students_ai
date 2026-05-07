package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrEmptyResponse = errors.New("empty response from api")
	ErrEmptyPrompt   = errors.New("prompt cannot be empty")
)

type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type request struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Stream   bool      `json:"stream,omitempty"`
}

type response struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

type streamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *Client) GetChatResponse(ctx context.Context, prompt string) (string, error) {
	if prompt == "" {
		return "", ErrEmptyPrompt
	}

	req := request{
		Model: "gpt-4o",
		Messages: []message{
			{
				Role:    "system",
				Content: "Sen CHDPU universiteti uchun maxsus yaratilgan AI yordamchisan. Sening vazifang faqat universitet, darslar, davomat va akademik monitoring bo'yicha savollarga javob berish. Agar foydalanuvchi universitetga aloqasi bo'lmagan (masalan, geografiya, umumiy bilim, ovqat pishirish va h.k.) savol bersa, muloyimlik bilan: 'Kechirasiz, men faqat universitet tizimi va akademik monitoring masalalari bo'yicha yordam berish uchun mo'ljallanganman' deb javob ber.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return "", err
	}
	defer body.Close()

	var resp response
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", ErrEmptyResponse
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *Client) GetChatResponseStream(ctx context.Context, prompt string, onChunk func(string)) error {
	if prompt == "" {
		return ErrEmptyPrompt
	}

	req := request{
		Model: "gpt-4o",
		Messages: []message{
			{
				Role:    "system",
				Content: "Sen CHDPU universiteti uchun maxsus yaratilgan AI yordamchisan. Sening vazifang faqat universitet, darslar, davomat va akademik monitoring bo'yicha savollarga javob berish. Agar foydalanuvchi universitetga aloqasi bo'lmagan (masalan, geografiya, umumiy bilim, ovqat pishirish va h.k.) savol bersa, muloyimlik bilan: 'Kechirasiz, men faqat universitet tizimi va akademik monitoring masalalari bo'yicha yordam berish uchun mo'ljallanganman' deb javob ber.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: true,
	}

	body, err := c.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer body.Close()

	return c.readStream(ctx, body, onChunk)
}

func (c *Client) doRequest(ctx context.Context, payload request) (io.ReadCloser, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("api error: status %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (c *Client) readStream(ctx context.Context, body io.Reader, onChunk func(string)) error {
	scanner := bufio.NewScanner(body)

	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		line := scanner.Text()

		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		if data == "[DONE]" {
			return nil
		}

		var chunk streamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}

		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			onChunk(chunk.Choices[0].Delta.Content)
		}
	}

	return scanner.Err()
}
