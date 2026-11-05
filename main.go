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
	redisClient, err := database.ConnectRedis(cfg.RedisURL)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	defer redisClient.Close()

	router := routes.Setup(db, redisClient, cfg)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
