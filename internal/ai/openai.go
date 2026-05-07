package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
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
				Content: "Sen CHDPU universiteti uchun maxsus yaratilgan AI yordamchisan. Sening vazifang faqat universitet, darslar, davomat va akademik monitoring bo'yicha savollarga javob berish. Agar foydalanuvchi universitetga aloqasi bo'lmagan (masalan, geografiya, umumiy bilim, ovqat pishirish va h.k.) savol bersa, muloyimlik bilan: 'Kechirasiz, men faqat universitet tizimi va akademik monitoring masalalari bo'yicha yordam berish uchun mo'ljallanganman' deb javob ber.",
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

func AudioToText(apiKey string, filePath string) (string, error) {
	url := "https://api.openai.com/v1/audio/transcriptions"

	// 1. Faylni ochish
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("faylni ochishda xato: %v", err)
	}
	defer file.Close()

	// 2. Multipart form yaratish
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return "", fmt.Errorf("form file yaratishda xato: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("faylni nusxalashda xato: %v", err)
	}

	_ = writer.WriteField("model", "whisper-1")
	//_ = writer.WriteField("language", "uz")
	writer.Close()

	// 3. Request yuborish
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("request yaratishda xato: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("so'rov yuborishda xato: %v", err)
	}
	defer resp.Body.Close()

	// Xatolikni tekshirish (StatusCode 200 bo'lmasa)
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API xatosi (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	// 4. Natijani DECODE qilish (Sizda shu qism tushib qolgan edi)
	var result struct {
		Text string `json:"text"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("JSON decode xatosi: %v", err)
	}

	return result.Text, nil
}
