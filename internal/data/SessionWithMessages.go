package data

// SessionWithMessages session va uning barcha messagelarini qaytaradi
type SessionWithMessages struct {
	Session  Session   `json:"session"`
	Messages []Message `json:"messages"`
}
