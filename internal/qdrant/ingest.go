package qdrant

import (
	"chdpu-ai-monitor/internal/ai"
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/qdrant/go-client/qdrant"
)

type StudentActivityDB struct {
	ID            uint64 `db:"id"`
	FullName      string `db:"full_name"`
	Group         string `db:"student_group"`
	AcademicYear  string `db:"academic_year"`
	Semester      string `db:"semester"`
	SubjectName   string `db:"subject_name"`
	ActivityValue string `db:"activity_value"`
	ActivityType  string `db:"activity_type"`
}

func ImportFromPostgres(db *sqlx.DB, qdrantClient *qdrant.Client, apiKey string) error {
	var activities []StudentActivityDB
	query := `
		SELECT id, full_name, student_group, academic_year, semester, subject_name, activity_value, activity_type 
		FROM student_activities 
		WHERE is_embedded = FALSE
		`

	err := db.Select(&activities, query)
	if err != nil {
		return fmt.Errorf("postgres'dan ma'lumot olishda xato: %v", err)
	}

	if len(activities) == 0 {
		log.Println("[INFO] Hamma ma'lumotlar allaqachon embedding qilingan.")
		return nil
	}

	log.Printf("[INFO] %d ta yangi qator topildi. Embedding boshlanmoqda...\n", len(activities))

	for _, act := range activities {
		// 2. AI uchun boyitilgan matnli kontekst tayyorlash
		contextText := fmt.Sprintf(
			"Talaba: %s, Guruhi: %s, O'quv yili: %s, Semestr: %s, Fan: %s, Turi: %s, Natija: %s",
			act.FullName, act.Group, act.AcademicYear, act.Semester, act.SubjectName, act.ActivityType, act.ActivityValue,
		)

		// 3. OpenAI'dan vektor olish
		vector, err := ai.GetEmbedding(context.Background(), apiKey, contextText)
		if err != nil {
			log.Printf("[XATO] Embedding olishda xato (ID: %d): %v\n", act.ID, err)
			continue
		}

		// 4. Qdrant nuqtasini (Point) shakllantirish
		point := &qdrant.PointStruct{
			Id:      qdrant.NewIDNum(act.ID),
			Vectors: qdrant.NewVectors(vector...),
			Payload: qdrant.NewValueMap(map[string]interface{}{
				"full_name":      act.FullName,
				"group":          act.Group,
				"academic_year":  act.AcademicYear,
				"semester":       act.Semester,
				"subject_name":   act.SubjectName,
				"activity_type":  act.ActivityType,
				"activity_value": act.ActivityValue,
				"text_context":   contextText, // Qidiruv vaqtida javob generatsiya qilish uchun kerak
			}),
		}

		_, err = qdrantClient.Upsert(context.Background(), &qdrant.UpsertPoints{
			CollectionName: "students_llm",
			Points:         []*qdrant.PointStruct{point},
		})
		if err != nil {
			log.Printf("[XATO] Qdrant'ga saqlashda xato (ID: %d): %v\n", act.ID, err)
			continue
		}

		_, err = db.Exec("UPDATE student_activities SET is_embedded = TRUE, updated_at = NOW() WHERE id = $1", act.ID)
		if err != nil {
			log.Printf("[XATO] Postgres statusini yangilashda xato (ID: %d): %v\n", act.ID, err)
		} else {
			log.Printf("[OK] ID: %d muvaffaqiyatli embedding qilindi.\n", act.ID)
		}
	}

	return nil
}
