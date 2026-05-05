package cmd

import (
	"chdpu-ai-monitor/internal/qdrant"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var ingestCmd = &cobra.Command{
	Use:   "make:ingest",
	Short: "Postgres'dan yangi ma'lumotlarni olib Qdrant'ga embedding qilish",
	Run: func(cmd *cobra.Command, args []string) {

		logDir := "storage/logs"
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatalf("Log papkasini yaratib bo'lmadi: %v", err)
		}

		logFileName := fmt.Sprintf("ingest-%s.log", time.Now().Format("2006-01-02"))
		logFilePath := filepath.Join(logDir, logFileName)

		logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Log faylini ochib bo'lmadi: %v", err)
		}
		defer logFile.Close()

		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

		log.Println("--- Ingest (Postgres -> Qdrant) jarayoni boshlandi ---")

		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Println("[XATO] OPENAI_API_KEY topilmadi")
			return
		}

		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			dbURL = "postgres://user:password@localhost:5433/student_llm?sslmode=disable"
		}

		db, err := sqlx.Connect("postgres", dbURL)
		if err != nil {
			log.Printf("[XATO] Postgres ulanish xatosi: %v\n", err)
			return
		}
		defer db.Close()

		client, err := qdrant.InitClient()
		if err != nil {
			log.Printf("[XATO] Qdrant ulanish xatosi: %v\n", err)
			return
		}
		defer client.Close()

		log.Println("Embedding jarayoni boshlandi...")

		err = qdrant.ImportFromPostgres(db, client, apiKey)
		if err != nil {
			log.Printf("[XATO] Ingest jarayonida xatolik: %v\n", err)
		} else {
			log.Println("--- Ingest jarayoni muvaffaqiyatli yakunlandi ---")
		}
	},
}
