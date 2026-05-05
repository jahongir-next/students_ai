package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

// GetChatResponse OpenAI API orqali javob oladi
func GetChatResponse(ctx context.Context, apiKey string, prompt string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	reqBody := ChatRequest{
		Model: "gpt-4o", // yoki "gpt-3.5-turbo"
		Messages: []ChatMessage{
			{
				Role:    "system",
				Content: "Sen talabalar monitoringi bo'yicha yordamchi AI assistentsan. Berilgan kontekstdan tashqariga chiqma.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// 2. HTTP so'rov yaratish
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 3. So'rovni yuborish
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 4. Xatolikni tekshirish
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI xatosi: status %d", resp.StatusCode)
	}

	// 5. Javobni o'qish
	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) > 0 {
		return chatResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("OpenAI'dan bo'sh javob qaytdi")
}
