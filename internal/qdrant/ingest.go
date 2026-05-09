package qdrant

import (
	"chdpu-ai-monitor/internal/ai"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
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

type TeacherLogDB struct {
	ID                 int       `db:"id"`
	TeacherName        string    `db:"teacher_name"`
	GroupName          string    `db:"group_name"`
	SubjectName        string    `db:"subject_name"`
	LessonType         string    `db:"lesson_type"`
	AcademicYear       string    `db:"academic_year"`
	Semester           string    `db:"semester"`
	ScheduledStartTime string    `db:"scheduled_start_time"`
	ActualArrivalTime  time.Time `db:"actual_arrival_time"`
	ArrivalDate        time.Time `db:"arrival_date"`
	Status             string    `db:"status"`
	DelayMinutes       int       `db:"delay_minutes"`
}

func ImportTeachersToQdrant(db *sqlx.DB, qdrantClient *qdrant.Client, apiKey string) error {
	var logs []TeacherLogDB
	query := `
       SELECT id, teacher_name, group_name, subject_name, lesson_type, academic_year, semester, 
              scheduled_start_time, actual_arrival_time, arrival_date, status, delay_minutes 
       FROM teacher_attendance_logs 
       WHERE is_embedded = FALSE
       `

	err := db.Select(&logs, query)
	if err != nil {
		return fmt.Errorf("postgres'dan ma'lumot olishda xato: %v", err)
	}

	if len(logs) == 0 {
		log.Println("[INFO] Hamma o'qituvchi ma'lumotlari allaqachon embedding qilingan.")
		return nil
	}

	log.Printf("[INFO] %d ta yangi o'qituvchi qatori topildi. Embedding boshlanmoqda...\n", len(logs))

	for _, act := range logs {
		contextText := fmt.Sprintf(
			"O'qituvchi: %s, O'quv yili: %s, Semestr: %s, Guruhi: %s, Fan: %s, Dars turi: %s. Kuni: %s, Rejadagi vaqt: %s, Kelgan vaqti: %s. Holati: %s, Kechikish: %d daqiqa.",
			act.TeacherName, act.AcademicYear, act.Semester, act.GroupName, act.SubjectName, act.LessonType,
			act.ArrivalDate.Format("2006-01-02"), act.ScheduledStartTime, act.ActualArrivalTime.Format("15:04"),
			act.Status, act.DelayMinutes,
		)

		// 2. OpenAI'dan vektor olish
		vector, err := ai.GetEmbedding(context.Background(), apiKey, contextText)
		if err != nil {
			log.Printf("[XATO] O'qituvchi ID %d uchun embedding olishda xato: %v\n", act.ID, err)
			continue
		}

		qdrantPointID := uuid.New().String()

		point := &qdrant.PointStruct{
			Id:      qdrant.NewIDUUID(qdrantPointID),
			Vectors: qdrant.NewVectors(vector...),
			Payload: qdrant.NewValueMap(map[string]interface{}{
				"entity_type":          "teacher",
				"db_id":                act.ID,
				"teacher_name":         act.TeacherName,
				"group_name":           act.GroupName,
				"subject_name":         act.SubjectName,
				"lesson_type":          act.LessonType,
				"academic_year":        act.AcademicYear,
				"semester":             act.Semester,
				"scheduled_start_time": act.ScheduledStartTime,
				"actual_arrival_time":  act.ActualArrivalTime.Format("2006-01-02T15:04:05Z"),
				"arrival_date":         act.ArrivalDate.Format("2006-01-02"),
				"status":               act.Status,
				"delay_minutes":        act.DelayMinutes,
				"text_context":         contextText,
			}),
		}

		// 4. Qdrant'ga saqlash (umumiy students_llm kolleksiyasiga)
		_, err = qdrantClient.Upsert(context.Background(), &qdrant.UpsertPoints{
			CollectionName: "students_llm",
			Points:         []*qdrant.PointStruct{point},
		})
		if err != nil {
			log.Printf("[XATO] Qdrant'ga saqlashda xato (O'qituvchi ID: %d): %v\n", act.ID, err)
			continue
		}

		_, err = db.Exec("UPDATE teacher_attendance_logs SET is_embedded = TRUE WHERE id = $1", act.ID)
		if err != nil {
			log.Printf("[XATO] Postgres statusini yangilashda xato (O'qituvchi ID: %d): %v\n", act.ID, err)
		} else {
			log.Printf("[OK] O'qituvchi ID: %d muvaffaqiyatli embedding qilindi.\n", act.ID)
		}
	}

	return nil
}
