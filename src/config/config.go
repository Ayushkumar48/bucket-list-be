package config

import "os"

type Config struct {
	DatabaseURL string
	Port        string
	RedisURL    string

	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	MagicLinkURL string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		Port:         port,
		RedisURL:     getenv("REDIS_URL", "redis://localhost:6379/0"),
		SMTPHost:     getenv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     getenv("SMTP_PORT", "587"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		MagicLinkURL: getenv("MAGIC_LINK_URL", "http://localhost:3000/auth/verify"),
	}
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
