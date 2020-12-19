package repo

import (
	"api/entities"
	"database/sql"
)

// ChatRepo ...
type ChatRepo struct {
	db *sql.DB
}

// NewChatRepo ...
func NewChatRepo(db *sql.DB) *ChatRepo {
	return &ChatRepo{
		db: db,
	}
}

// Create ...
func (c *ChatRepo) Create(chat *entities.Chat) error {
	return c.db.QueryRow(`INSERT INTO chats (name) VALUES ($1) RETURNING id`, chat.Name).Scan(&chat.ID)
}

// GetAll ...
func (c *ChatRepo) GetAll() ([]*entities.Chat, error) {
	rows, err := c.db.Query(`SELECT * FROM chats`)
	if err != nil {
		return nil, err
	}

	var chats []*entities.Chat
	for rows.Next() {
		var (
			id   int
			name string
		)

		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}

		chat := &entities.Chat{
			ID:   id,
			Name: name,
		}

		chats = append(chats, chat)
	}
	return chats, nil
}
