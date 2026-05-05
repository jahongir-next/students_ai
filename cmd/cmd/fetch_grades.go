package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

type StudentData struct {
	FullName     string            `json:"full_name"`
	PaymentType  string            `json:"payment_type"`
	AcademicYear string            `json:"academic_year"`
	Semester     string            `json:"semester"`
	Group        string            `json:"group"`
	Grades       map[string]string `json:"grades"`
}

var fetchGradesCmd = &cobra.Command{
	Use:   "make:fetch-grades",
	Short: "Google Sheets-dan baholarni olib Postgres-ga yozish",
	Run: func(cmd *cobra.Command, args []string) {

		logDir := "storage/logs"
		os.MkdirAll(logDir, 0755)
		logFileName := fmt.Sprintf("fetch-%s.log", time.Now().Format("2006-01-02"))
		logFilePath := filepath.Join(logDir, logFileName)

		logFile, _ := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer logFile.Close()

		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

		log.Println("--- Fetch jarayoni boshlandi ---")

		// 2. Bazaga ulanish
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			log.Println("[XATO] DATABASE_URL .env faylda topilmadi")
			return
		}

		db, err := sqlx.Connect("postgres", dbURL)
		if err != nil {
			log.Printf("[XATO] Postgres ulanish xatosi: %v\n", err)
			return
		}
		defer db.Close()

		urls := []string{
			os.Getenv("FIRST_GRADE_URL"),
			os.Getenv("SECOND_GRADE_URL"),
		}

		var wg sync.WaitGroup
		for _, url := range urls {
			wg.Add(1)
			go func(targetURL string) {
				defer wg.Done()
				processFetch(db, targetURL)
			}(url)
		}
		wg.Wait()

		log.Println("--- Fetch jarayoni muvaffaqiyatli yakunlandi ---")
	},
}

func init() {
	rootCmd.AddCommand(fetchGradesCmd)
}

func processFetch(db *sqlx.DB, url string) {
	log.Printf("[INFO] Ma'lumot olinmoqda: %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[XATO] URL so'rovda xato: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var students []StudentData
	if err := json.Unmarshal(body, &students); err != nil {
		log.Printf("[XATO] JSON parsing xatosi: %v\n", err)
		return
	}

	for _, s := range students {
		for subject, grade := range s.Grades {
			// Universal migrationga mos INSERT + ON CONFLICT (Duplicate oldini olish)
			query := `
				INSERT INTO student_activities (
					full_name, student_group, academic_year, semester, payment_type, 
					activity_type, subject_name, activity_value, is_grade
				) VALUES ($1, $2, $3, $4, $5, 'grade', $6, $7, true)
				ON CONFLICT (full_name, academic_year, semester, subject_name) 
				WHERE (activity_type = 'grade')
				DO UPDATE SET 
					activity_value = EXCLUDED.activity_value,
					updated_at = CURRENT_TIMESTAMP;`

			_, err := db.Exec(query,
				s.FullName, s.Group, s.AcademicYear, s.Semester, s.PaymentType,
				subject, grade,
			)
			if err != nil {
				log.Printf("[XATO] Bazaga yozishda xato (%s): %v\n", s.FullName, err)
			}
		}
	}
}
