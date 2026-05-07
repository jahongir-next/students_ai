package main

import (
	openai "chdpu-ai-monitor/cmd/api"
	"chdpu-ai-monitor/cmd/cmd"
	"chdpu-ai-monitor/internal/api"
	"chdpu-ai-monitor/internal/qdrant"
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	godotenv.Load()

	cmd.Execute()

	// Database connection
	db, err := sql.Open(
		"postgres",
		"postgresql://user:password@localhost:5433/student_llm?sslmode=disable",
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	qClient, err := qdrant.InitClient()
	if err != nil {
		log.Fatal("Qdrant ulanish xatosi:", err)
	}
	defer qClient.Close()

	h := &api.Handler{
		OpenAI:      client,
		Logger:      logger,
		Qdrant:      qClient,
		SessionRepo: openai.NewSessionRepository(db),
	}

	router := httprouter.New()

	router.POST("/api/ask", api.AskHandler(qClient))
	router.POST("/api/ask-voice", api.VoiceAskHandler(qClient))
	router.POST("/api/ask/stream", h.ChatStreamHandler)

	// SESSION ENDPOINTS
	router.POST("/api/sessions", h.CreateSessionHandler)
	router.GET("/api/sessions/:id", h.GetSessionHandler)
	router.GET("/api/sessions", h.ListSessionsHandler)
	router.DELETE("/api/sessions/:id", h.DeleteSessionHandler)

	handler := openai.EnableCORS(router)
	log.Println("Server 8085 portida ishlamoqda...")
	log.Fatal(http.ListenAndServe(":8085", handler))
}
