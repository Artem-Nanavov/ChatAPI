package main

import (
	"api/database"
	"api/server"
	"api/utils"
)

func main() {
	// Database
	db := database.NewDatabase(&database.Config{
		DSN: utils.GetEnv(
			"DSN", "user=postgres password=1234 dbname=postgres port=5432 host=localhost sslmode=disable"),
	})
	if err := db.Open(); err != nil {
		return
	}
	if err := db.Migrate("./database/migrations"); err != nil {
		return
	}
	defer db.Close()

	// Server
	serv := server.NewServer(db, &server.Config{
		Port:      utils.GetEnv("PORT", "8080"),
		JWTSecret: utils.GetEnv("JWTSECRET", "asdasdasdasd"),
		Salt:      utils.GetEnv("SALT", "asdasdasdasd"),
	})

	if err := serv.Run(); err != nil {
		return
	}

}
