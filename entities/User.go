package entities

// User ...
type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	IsOnline bool `json:"is_online"`
}
