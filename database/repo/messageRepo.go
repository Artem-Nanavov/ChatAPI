package repo

import (
	"api/entities"
	"database/sql"
)

// MessageRepo ...
type MessageRepo struct {
	db *sql.DB
}

// NewMessageRepo ...
func NewMessageRepo(db *sql.DB) *MessageRepo {
	return &MessageRepo{
		db: db,
	}
}

// Create ...
func (m *MessageRepo) Create(message *entities.Message) error {
	return m.db.QueryRow(`INSERT INTO messages (text, owner_id, chat_id) VALUES ($1, $2, $3) RETURNING id`,
		message.Text, message.OwnerID, message.ChatID).Scan(&message.ID)
}

// GetAllByChatID ...
func (m *MessageRepo) GetAllByChatID(id int) ([]*entities.Message, error) {
	rows, err := m.db.Query(`SELECT * FROM messages WHERE chat_id=$1`, id)
	if err != nil {
		return nil, err
	}

	var messages []*entities.Message
	for rows.Next() {
		var (
			id      int
			text    string
			ownerID int
			chatID  int
		)

		if err := rows.Scan(&id, &text, &ownerID, &chatID); err != nil {
			return nil, err
		}

		msg := &entities.Message{
			ID:      id,
			Text:    text,
			OwnerID: ownerID,
			ChatID:  chatID,
		}

		messages = append(messages, msg)
	}
	return messages, nil
}
