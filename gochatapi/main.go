package main

import (
	"api/chat"
	"api/database"
	"api/server"
	"api/utils"
)

func main() {
	// Database
	db := database.NewDatabase(&database.Config{
		DSN: utils.GetEnv("DSN", "asdasd"),
	})
	if err := db.Open(); err != nil {
		return
	}
	if err := db.Migrate("./database/migrations"); err != nil {
		return
	}
	defer db.Close()

	// Chat
	wschat := chat.NewWebSocketChat(db, &chat.Config{
		Port: utils.GetEnv("WEBSOCKET_PORT", "8000"),
	})

	// Server
	serv := server.NewServer(db, &server.Config{
		Port:      utils.GetEnv("PORT", "8080"),
		JWTSecret: utils.GetEnv("JWTSECRET", "random"),
		Salt:      utils.GetEnv("SALT", "random"),
	})

	go func() {
		if err := wschat.Run(); err != nil {
			return
		}
	}()

	if err := serv.Run(); err != nil {
		return
	}

}
