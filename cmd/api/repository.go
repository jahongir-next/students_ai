package openai

import (
	"chdpu-ai-monitor/internal/data"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// CreateSession yangi session yaratadi
func (r *SessionRepository) CreateSession(ctx context.Context, userID, title string) (*data.Session, error) {
	session := &data.Session{
		ID:        uuid.New().String(),
		UserID:    userID,
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO sessions (id, user_id, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, title, created_at, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		session.ID, session.UserID, session.Title, session.CreatedAt, session.UpdatedAt,
	).Scan(&session.ID, &session.UserID, &session.Title, &session.CreatedAt, &session.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession sessionni ID bo'yicha oladi
func (r *SessionRepository) GetSession(ctx context.Context, sessionID string) (*data.Session, error) {
	session := &data.Session{}
	query := `SELECT id, user_id, title, created_at, updated_at FROM sessions WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID, &session.UserID, &session.Title, &session.CreatedAt, &session.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return session, nil
}

// ListSessions user uchun barcha sessionlarni qaytaradi
func (r *SessionRepository) ListSessions(ctx context.Context, userID string, limit int) ([]data.Session, error) {
	query := `
		SELECT id, user_id, title, created_at, updated_at 
		FROM sessions 
		WHERE user_id = $1 
		ORDER BY updated_at DESC 
		LIMIT $2
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []data.Session
	for rows.Next() {
		var s data.Session
		if err := rows.Scan(&s.ID, &s.UserID, &s.Title, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}

	return sessions, nil
}

// UpdateSessionTitle session titleni yangilaydi
func (r *SessionRepository) UpdateSessionTitle(ctx context.Context, sessionID, title string) error {
	query := `UPDATE sessions SET title = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, title, sessionID)
	return err
}

// DeleteSession sessionni o'chiradi (messages ham cascade o'chiriladi)
func (r *SessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, sessionID)
	return err
}

// AddMessage sessionga yangi message qo'shadi
func (r *SessionRepository) AddMessage(ctx context.Context, sessionID, role, content string) (*data.Message, error) {
	message := &data.Message{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Role:      role,
		Content:   content,
		CreatedAt: time.Now(),
	}

	query := `
		INSERT INTO messages (id, session_id, role, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, session_id, role, content, created_at
	`

	err := r.db.QueryRowContext(ctx, query,
		message.ID, message.SessionID, message.Role, message.Content, message.CreatedAt,
	).Scan(&message.ID, &message.SessionID, &message.Role, &message.Content, &message.CreatedAt)

	if err != nil {
		return nil, err
	}

	// Session updated_at ni yangilash
	_, _ = r.db.ExecContext(ctx, `UPDATE sessions SET updated_at = NOW() WHERE id = $1`, sessionID)

	return message, nil
}

// GetMessages session uchun barcha messagelarni qaytaradi
func (r *SessionRepository) GetMessages(ctx context.Context, sessionID string) ([]data.Message, error) {
	query := `
		SELECT id, session_id, role, content, created_at 
		FROM messages 
		WHERE session_id = $1 
		ORDER BY created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []data.Message
	for rows.Next() {
		var m data.Message
		if err := rows.Scan(&m.ID, &m.SessionID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

// GetSessionWithMessages session va uning messagelarini qaytaradi
func (r *SessionRepository) GetSessionWithMessages(ctx context.Context, sessionID string) (*data.SessionWithMessages, error) {
	session, err := r.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	messages, err := r.GetMessages(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &data.SessionWithMessages{
		Session:  *session,
		Messages: messages,
	}, nil
}
