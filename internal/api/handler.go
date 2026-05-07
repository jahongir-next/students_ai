package api

import (
	openai "chdpu-ai-monitor/cmd/api"
	"chdpu-ai-monitor/internal/ai"
	"context"
	_ "context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/qdrant/go-client/qdrant"
)

type Handler struct {
	OpenAI      *openai.Client
	Logger      *log.Logger
	Qdrant      *qdrant.Client
	SessionRepo *openai.SessionRepository
}

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
			Limit:          qdrant.PtrOf(uint64(30)),
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
			Limit:          qdrant.PtrOf(uint64(30)),
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

// CreateSessionHandler yangi session yaratadi
func (h *Handler) CreateSessionHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		UserID string `json:"user_id"` // optional
		Title  string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		input.Title = "New Conversation"
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	session, err := h.SessionRepo.CreateSession(ctx, input.UserID, input.Title)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// GetSessionHandler sessionni oladi
func (h *Handler) GetSessionHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	sessionID := params.ByName("id")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	sessionWithMessages, err := h.SessionRepo.GetSessionWithMessages(ctx, sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessionWithMessages)
}

// ListSessionsHandler user sessionlarini listlaydi
func (h *Handler) ListSessionsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	userID := r.URL.Query().Get("user_id")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	sessions, err := h.SessionRepo.ListSessions(ctx, userID, 50)
	if err != nil {
		http.Error(w, "Failed to list sessions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sessions)
}

// DeleteSessionHandler sessionni o'chiradi
func (h *Handler) DeleteSessionHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if r.Method != http.MethodDelete {
		w.Header().Set("Allow", http.MethodDelete)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := params.ByName("id")

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	if err := h.SessionRepo.DeleteSession(ctx, sessionID); err != nil {
		http.Error(w, "Failed to delete session", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// =========================================================
// 💬 YANGILANGAN CHAT STREAM (SESSION-AWARE)
// =========================================================

func (h *Handler) ChatStreamHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// --- validation ---
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		SessionID string `json:"session_id"` // MUHIM: session ID
		Prompt    string `json:"prompt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if input.Prompt == "" {
		http.Error(w, "Prompt is required", http.StatusBadRequest)
		return
	}

	if input.SessionID == "" {
		http.Error(w, "Session ID is required", http.StatusBadRequest)
		return
	}

	// --- STREAM HEADERS ---
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 90*time.Second)
	defer cancel()

	// Sessionni tekshirish
	session, err := h.SessionRepo.GetSession(ctx, input.SessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			sendStreamError(w, flusher, "Session not found")
			return
		}
		sendStreamError(w, flusher, "Failed to validate session")
		return
	}

	fmt.Fprintf(w, "event: start\ndata: {\"status\":\"processing\"}\n\n")
	flusher.Flush()

	// User messageini saqlash
	_, err = h.SessionRepo.AddMessage(ctx, session.ID, "user", input.Prompt)
	if err != nil {
		sendStreamError(w, flusher, "Failed to save user message")
		return
	}

	// =========================================================
	// 🔥 1. CONVERSATION HISTORY OLISH
	// =========================================================
	messages, err := h.SessionRepo.GetMessages(ctx, session.ID)
	if err != nil {
		sendStreamError(w, flusher, "Failed to load conversation history")
		return
	}

	// =========================================================
	// 🔥 2. EMBEDDING
	// =========================================================
	apiKey := os.Getenv("OPENAI_API_KEY")

	queryVector, err := ai.GetEmbedding(ctx, apiKey, input.Prompt)
	if err != nil {
		sendStreamError(w, flusher, "Embedding error")
		return
	}

	// =========================================================
	// 🔥 3. QDRANT SEARCH
	// =========================================================
	searchRes, err := h.Qdrant.Query(ctx, &qdrant.QueryPoints{
		CollectionName: "students_llm",
		Query:          qdrant.NewQuery(queryVector...),
		Limit:          qdrant.PtrOf(uint64(10)),
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		sendStreamError(w, flusher, "Qdrant search error")
		return
	}

	// =========================================================
	// 🔥 4. CONTEXT YIG'ISH
	// =========================================================
	var contextText strings.Builder

	for _, hit := range searchRes {
		if val, ok := hit.Payload["text_context"]; ok {
			if strVal := val.GetStringValue(); strVal != "" {
				contextText.WriteString(strVal)
				contextText.WriteString("\n")
			}
		}
	}

	// =========================================================
	// 🔥 5. CONVERSATION HISTORY FORMATLASH
	// =========================================================
	var conversationHistory strings.Builder

	// Oxirgi 5 ta message (context uchun, agar kerak bo'lsa)
	startIdx := 0
	if len(messages) > 5 {
		startIdx = len(messages) - 5
	}

	for _, msg := range messages[startIdx:] {
		role := "User"
		if msg.Role == "assistant" {
			role = "Assistant"
		}
		conversationHistory.WriteString(fmt.Sprintf("%s: %s\n", role, msg.Content))
	}

	// =========================================================
	// 🔥 6. FINAL PROMPT (RAG + HISTORY)
	// =========================================================
	finalPrompt := fmt.Sprintf(
		`Quyidagi ma'lumotlarga va oldingi suhbatga asoslanib javob ber:
 
KONTEKST (Ma'lumotlar bazasidan):
%s
 
OLDINGI SUHBAT:
%s
 
JORIY SAVOL: %s`,
		contextText.String(),
		conversationHistory.String(),
		input.Prompt,
	)

	fmt.Fprintf(w, "event: start\ndata: {\"status\":\"streaming\"}\n\n")
	flusher.Flush()

	// =========================================================
	// 🔥 7. STREAM OPENAI va JAVOBNI SAQLASH
	// =========================================================
	var fullResponse strings.Builder

	err = h.OpenAI.GetChatResponseStream(ctx, finalPrompt, func(chunk string) {
		fullResponse.WriteString(chunk)

		data := map[string]string{"content": chunk}
		jsonData, _ := json.Marshal(data)

		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
	})

	if err != nil {
		sendStreamError(w, flusher, "Stream error occurred")
		return
	}

	// Assistant javobini saqlash
	_, err = h.SessionRepo.AddMessage(ctx, session.ID, "assistant", fullResponse.String())
	if err != nil {
		// Log qilish, lekin streamga tasir qilmaslik
		fmt.Printf("Warning: failed to save assistant message: %v\n", err)
	}

	fmt.Fprintf(w, "event: done\ndata: {\"status\":\"completed\"}\n\n")
	flusher.Flush()
}

func sendStreamError(w http.ResponseWriter, flusher http.Flusher, msg string) {
	data := map[string]string{"error": msg}
	jsonData, _ := json.Marshal(data)

	fmt.Fprintf(w, "event: error\ndata: %s\n\n", jsonData)
	flusher.Flush()
}
