package qdrant

import (
	"os"
	"strconv"

	"github.com/qdrant/go-client/qdrant"
)

func InitClient() (*qdrant.Client, error) {
	
	qdrantHost := os.Getenv("QDRANT_HOST")
	qdrantPortStr := os.Getenv("QDRANT_PORT")

	port, err := strconv.Atoi(qdrantPortStr)
	if err != nil {
		port = 6334
	}

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: qdrantHost,
		Port: port,
	})

	return client, err
}
