package entities

// Message ...
type Message struct {
	ID      int    `json:"id"`
	Text    string `json:"text"`
	OwnerID int    `json:"owner_id"`
	ChatID int 
}
