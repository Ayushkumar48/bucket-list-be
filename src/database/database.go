package database

import (
	"bucket-list-be/src/models"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.User{}, &models.Session{}); err != nil {
		return nil, err
	}
	// The previous password-based schema made these columns mandatory. Keep
	// existing profile data, but allow a newly-created magic-link user to have
	// only an email. Password hashes are no longer used by this application.
	if err := db.Exec(`
		ALTER TABLE users ALTER COLUMN name DROP NOT NULL;
		ALTER TABLE users ALTER COLUMN username DROP NOT NULL;
		ALTER TABLE users DROP COLUMN IF EXISTS password_hash;
	`).Error; err != nil {
		return nil, err
	}
	return db, nil
}

// ConnectRedis creates and verifies a Redis-compatible client, including
// Dragonfly deployments that expose the Redis protocol.
func ConnectRedis(redisURL string) (*redis.Client, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(options)
	if err := client.Ping(context.Background()).Err(); err != nil {
		_ = client.Close()
		return nil, err
	}

	return client, nil
}
