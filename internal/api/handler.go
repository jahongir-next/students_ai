package api

import (
	"chdpu-ai-monitor/internal/ai"
	_ "context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/qdrant/go-client/qdrant"
)

type AskRequest struct {
	Question string `json:"question"`
}

type AskResponse struct {
	Answer string `json:"answer"`
}

func AskHandler(qClient *qdrant.Client) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var req AskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Noto'g'ri so'rov", http.StatusBadRequest)
			return
		}

		apiKey := os.Getenv("OPENAI_API_KEY")

		// 1. Savolni embedding qilish
		queryVector, err := ai.GetEmbedding(r.Context(), apiKey, req.Question)
		if err != nil {
			http.Error(w, "Embedding xatosi", http.StatusInternalServerError)
			return
		}

		// 2. Qdrant'dan o'xshash ma'lumotlarni qidirish (Search)
		searchRes, err := qClient.Query(r.Context(), &qdrant.QueryPoints{
			CollectionName: "students_llm",
			Query:          qdrant.NewQuery(queryVector...),
			Limit:          qdrant.PtrOf(uint64(10)),
			WithPayload:    qdrant.NewWithPayload(true),
		})
		if err != nil {
			http.Error(w, "Qdrant qidiruv xatosi", http.StatusInternalServerError)
			return
		}

		// 3. Topilgan ma'lumotlardan kontekst yig'ish
		contextText := ""
		for _, hit := range searchRes {
			// Payload'dan "text_context" maydonini olamiz
			if val, ok := hit.Payload["text_context"]; ok {
				// Qdrant Value turidan string qiymatni ajratib olish
				// GetStringValue() protobuff'dan keladigan metod
				if strVal := val.GetStringValue(); strVal != "" {
					contextText += strVal + "\n"
				}
			}
		}
		// internal/api/handler.go ichida
		//fmt.Println("OpenAI'ga ketayotgan kontekst:", contextText)

		// 4. OpenAI'ga kontekst bilan birga savolni yuborish
		// Bu yerda ai paketingizda ChatCompletion funksiyasi bor deb faraz qilamiz
		prompt := fmt.Sprintf(
			"Quyidagi ma'lumotlarga asoslanib savolga javob ber:\n\nKONTEKST:\n%s\n\nSAVOL: %s",
			contextText, req.Question,
		)

		answer, err := ai.GetChatResponse(r.Context(), apiKey, prompt)
		if err != nil {
			http.Error(w, "OpenAI javob xatosi", http.StatusInternalServerError)
			return
		}

		// 5. Javobni qaytarish
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AskResponse{Answer: answer})
	}
}

func VoiceAskHandler(qClient *qdrant.Client) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		apiKey := os.Getenv("OPENAI_API_KEY")

		// 1. Audio faylni qabul qilish (multipart/form-data)
		file, header, err := r.FormFile("audio")
		if err != nil {
			http.Error(w, "Audio fayl topilmadi", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 2. Vaqtinchalik fayl yaratish
		tempDir := "storage/temp"
		os.MkdirAll(tempDir, os.ModePerm)
		tempFile := filepath.Join(tempDir, header.Filename)

		f, err := os.Create(tempFile)
		if err != nil {
			http.Error(w, "Faylni saqlashda xato", http.StatusInternalServerError)
			return
		}
		io.Copy(f, file)
		f.Close()
		defer os.Remove(tempFile) // Ish bitgach o'chirish

		// 3. Whisper orqali matnga aylantirish (STT)
		// Eslatma: ai.AudioToText funksiyasini yuqorida gaplashganimizdek yaratib olishingiz kerak
		transcribedText, err := ai.AudioToText(apiKey, tempFile)
		fmt.Println(transcribedText)
		if err != nil {
			http.Error(w, "Ovozni matnga aylantirishda xato: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 4. Endi matn tayyor, uni mavjud RAG jarayoniga yo'naltiramiz
		// Bu qismi sizning AskHandler kodingiz bilan deyarli bir xil bo'ladi

		// a. Embed qilish (Savolni vektorga aylantirish)
		queryVector, err := ai.GetEmbedding(r.Context(), apiKey, transcribedText)
		if err != nil {
			http.Error(w, "Embedding xatosi", http.StatusInternalServerError)
			return
		}

		// b. Qdrantdan qidirish
		searchRes, err := qClient.Query(r.Context(), &qdrant.QueryPoints{
			CollectionName: "students_llm",
			Query:          qdrant.NewQuery(queryVector...),
			Limit:          qdrant.PtrOf(uint64(10)),
			WithPayload:    qdrant.NewWithPayload(true),
		})

		// 3. Topilgan ma'lumotlardan kontekst yig'ish
		contextText := ""
		for _, hit := range searchRes {
			if val, ok := hit.Payload["text_context"]; ok {
				if strVal := val.GetStringValue(); strVal != "" {
					contextText += strVal + "\n"
				}
			}
		}

		prompt := fmt.Sprintf("Kontekst:\n%s\n\nSavol: %s", contextText, transcribedText)
		answer, err := ai.GetChatResponse(r.Context(), apiKey, prompt)

		// 5. Natijani qaytarish
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"question": "%s", "answer": "%s"}`, transcribedText, answer)
	}
}
