package handlers

import (
	"bucket-list-be/src/config"
	"bucket-list-be/src/models"
	"bucket-list-be/src/utils"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	sessionDuration  = 30 * 24 * time.Hour
	magicLinkTTL     = 10 * time.Minute
	magicLinkKeyPref = "bucket-list:magic-link:"
)

type LoginRequest struct {
	Email string `json:"email"`
}

type VerifyMagicLinkRequest struct {
	Token string `json:"token"`
}

func RequestMagicLink(db *gorm.DB, redisClient *redis.Client, cfg config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, ok := utils.BindJSON[LoginRequest](c)
		if !ok {
			return
		}
		email := strings.ToLower(strings.TrimSpace(req.Email))
		if email == "" || !strings.Contains(email, "@") {
			utils.BadRequest(c, "A valid email is required")
			return
		}

		ctx := c.Request.Context()

		rawToken, err := randomToken()
		if err != nil {
			utils.InternalError(c)
			return
		}
		key := magicLinkKey(rawToken)
		if err := redisClient.Set(ctx, key, email, magicLinkTTL).Err(); err != nil {
			utils.InternalError(c)
			return
		}

		link := cfg.MagicLinkURL + "?token=" + url.QueryEscape(rawToken)
		sender := EmailSender{Host: cfg.SMTPHost, Port: cfg.SMTPPort, Username: cfg.SMTPUsername, Password: cfg.SMTPPassword}
		if err := sender.SendMagicLink(email, link); err != nil {
			_ = redisClient.Del(ctx, key).Err()
			utils.InternalError(c)
			return
		}
		utils.ResponseOk(c, gin.H{"message": "If that email can be used, a sign-in link has been sent"})
	}
}

func VerifyMagicLink(db *gorm.DB, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, ok := utils.BindJSON[VerifyMagicLinkRequest](c)
		if !ok {
			return
		}
		token := strings.TrimSpace(req.Token)
		if token == "" {
			utils.Unauthorized(c, "Invalid or expired magic link")
			return
		}
		ctx := c.Request.Context()
		// GETDEL atomically reads and deletes the token, making it single-use.
		email, err := redisClient.GetDel(ctx, magicLinkKey(token)).Result()
		if errors.Is(err, redis.Nil) {
			utils.Unauthorized(c, "Invalid or expired magic link")
			return
		}
		if err != nil {
			utils.InternalError(c)
			return
		}

		user, err := gorm.G[models.User](db).Where("email = ?", email).First(ctx)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = models.User{Email: email}
			if err := gorm.G[models.User](db).Create(ctx, &user); err != nil {
				// Another verification request may have created the same account.
				if errors.Is(err, gorm.ErrDuplicatedKey) {
					user, err = gorm.G[models.User](db).Where("email = ?", email).First(ctx)
				}
				if err != nil {
					utils.InternalError(c)
					return
				}
			}
		} else if err != nil {
			utils.InternalError(c)
			return
		}

		sessionToken, err := randomToken()
		if err != nil {
			utils.InternalError(c)
			return
		}
		session := models.Session{
			UserID: user.ID, TokenHash: hashToken(sessionToken), ClientType: "web",
			ExpiresAt: time.Now().UTC().Add(sessionDuration),
		}
		if err := gorm.G[models.Session](db).Create(ctx, &session); err != nil {
			utils.InternalError(c)
			return
		}
		utils.ResponseOk(c, gin.H{"message": "Login successful", "token": sessionToken, "token_type": "Bearer", "expires_at": session.ExpiresAt})
	}
}

func Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			utils.Unauthorized(c, "Missing or invalid authorization token")
			return
		}
		token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
		if token == "" {
			utils.Unauthorized(c, "Missing or invalid authorization token")
			return
		}
		result, err := gorm.G[models.Session](db).Where("token_hash = ?", hashToken(token)).Delete(c.Request.Context())
		if err != nil {
			utils.InternalError(c)
			return
		}
		if result == 0 {
			utils.Unauthorized(c, "Invalid or expired session")
			return
		}
		utils.ResponseOk(c, gin.H{"message": "Logout successful"})
	}
}

func Me(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			utils.Unauthorized(c)
			return
		}
		id, ok := userID.(uuid.UUID)
		if !ok {
			utils.InternalError(c)
			return
		}
		user, err := gorm.G[models.User](db).Where("id = ?", id).First(c.Request.Context())
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.Unauthorized(c)
			} else {
				utils.InternalError(c)
			}
			return
		}
		utils.ResponseOk(c, user)
	}
}

func randomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func magicLinkKey(token string) string {
	return magicLinkKeyPref + hashToken(token)
}
