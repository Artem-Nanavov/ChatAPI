package entities

import "time"

// Message ...
type Message struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	OwnerID   int       `json:"owner_id"`
	ChatID    int       `json:"chat_id"`
	CreatedAt time.Time `json:"created_at"`
}
