package cmd

import (
	"chdpu-ai-monitor/internal/qdrant" // O'zingizning loyiha nomingizga qarab o'zgartiring
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

var ingestTeachersCmd = &cobra.Command{
	Use:   "make:ingest-teachers",
	Short: "O'qituvchilar davomati ma'lumotlarini yagona students_llm kolleksiyasiga embedding qilish",
	Run: func(cmd *cobra.Command, args []string) {

		// 1. Loglarni sozlash
		logDir := "storage/logs"
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatalf("Log papkasini yaratib bo'lmadi: %v", err)
		}

		logFileName := fmt.Sprintf("ingest-teachers-%s.log", time.Now().Format("2006-01-02"))
		logFilePath := filepath.Join(logDir, logFileName)

		logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Log faylini ochib bo'lmadi: %v", err)
		}
		defer logFile.Close()

		// Ham terminalga, ham faylga yozish
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

		log.Println("--- O'qituvchilar Ingest (Postgres -> Qdrant) jarayoni boshlandi ---")

		// 2. Muhit o'zgaruvchilarini tekshirish
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			log.Println("[XATO] OPENAI_API_KEY topilmadi (.env faylni tekshiring)")
			return
		}

		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			// MacBook uchun default port 5433 (docker-compose faylingizdagi)
			dbURL = "postgres://user:password@localhost:5433/student_llm?sslmode=disable"
		}

		// 3. Ulanishlarni o'rnatish
		db, err := sqlx.Connect("postgres", dbURL)
		if err != nil {
			log.Printf("[XATO] Postgres ulanish xatosi: %v\n", err)
			return
		}
		defer db.Close()

		qdrantClient, err := qdrant.InitClient()
		if err != nil {
			log.Printf("[XATO] Qdrant ulanish xatosi: %v\n", err)
			return
		}
		defer qdrantClient.Close()

		log.Println("Baza tekshirilmoqda va embedding jarayoni boshlanmoqda...")

		// 4. Asosiy mantiqni chaqirish (Avvalgi turnlarda yozgan funksiyamiz)
		err = qdrant.ImportTeachersToQdrant(db, qdrantClient, apiKey)
		if err != nil {
			log.Printf("[XATO] Ingest jarayonida xatolik yuz berdi: %v\n", err)
		} else {
			log.Println("--- O'qituvchilar Ingest jarayoni muvaffaqiyatli yakunlandi ---")
		}
	},
}

func init() {
	// Buyruqni root command'ga ro'yxatdan o'tkazamiz
	rootCmd.AddCommand(ingestTeachersCmd)
}
