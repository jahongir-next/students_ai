package main

import (
	openai "chdpu-ai-monitor/cmd/api"
	"chdpu-ai-monitor/cmd/cmd"
	"chdpu-ai-monitor/internal/api"
	"chdpu-ai-monitor/internal/qdrant"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	godotenv.Load()

	cmd.Execute()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	qClient, err := qdrant.InitClient()
	if err != nil {
		log.Fatal("Qdrant ulanish xatosi:", err)
	}
	defer qClient.Close()

	h := &api.Handler{
		OpenAI: client,
		Logger: logger,
		Qdrant: qClient,
	}

	router := httprouter.New()

	router.POST("/api/ask", api.AskHandler(qClient))
	router.POST("/api/ask-voice", api.VoiceAskHandler(qClient))
	router.POST("/api/ask/stream", h.ChatStreamHandler)

	log.Println("Server 8085 portida ishlamoqda...")
	log.Fatal(http.ListenAndServe(":8085", router))
}
