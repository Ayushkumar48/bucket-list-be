package main

import (
	"bucket-list-be/src/config"
	"bucket-list-be/src/database"
	"bucket-list-be/src/routes"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	router := routes.Setup(db, cfg)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
