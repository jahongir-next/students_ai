package main

import (
	"chdpu-ai-monitor/cmd/cmd"
	"chdpu-ai-monitor/internal/api"
	"chdpu-ai-monitor/internal/qdrant"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {
	godotenv.Load()

	cmd.Execute()

	qClient, err := qdrant.InitClient()
	if err != nil {
		log.Fatal("Qdrant ulanish xatosi:", err)
	}
	defer qClient.Close()

	router := httprouter.New()

	router.POST("/api/ask", api.AskHandler(qClient))

	log.Println("Server 8085 portida ishlamoqda...")
	log.Fatal(http.ListenAndServe(":8085", router))
}
