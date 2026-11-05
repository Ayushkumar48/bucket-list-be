package handlers

import (
	"bucket-list-be/src/models"
	"bucket-list-be/src/utils/schema"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"bucket-list-be/src/utils"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const sessionDuration = 30 * 24 * time.Hour

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

type SignupRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Signup(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		req, ok := utils.BindJSON[SignupRequest](c)
		if !ok {
			return
		}

		issues := schema.SignupSchema.Validate(&req)
		if !utils.ValidateZog(c, issues) {
			return
		}

		// Check email.
		_, err := gorm.G[models.User](db).
			Where("email = ?", req.Email).
			First(ctx)

		if err == nil {
			utils.ConflictRequest(c, "An account with this email already exists")
			return
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.InternalError(c)
			return
		}

		// Check username.
		_, err = gorm.G[models.User](db).
			Where("username = ?", req.Username).
			First(ctx)

		if err == nil {
			utils.ConflictRequest(c, "Username already exists")
			return
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.InternalError(c)
			return
		}

		// Hash password using Argon2id.
		passwordHash, err := argon2id.CreateHash(
			req.Password,
			argon2id.DefaultParams,
		)
		if err != nil {
			utils.InternalError(c)
			return
		}

		user := models.User{
			Name:         req.Name,
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: passwordHash,
		}

		// Create user.
		err = gorm.G[models.User](db).Create(ctx, &user)
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				c.JSON(http.StatusConflict, gin.H{
					"error": "Username or email already exists",
				})
				return
			}

			utils.InternalError(c)
			return
		}

		utils.ResponseOk(c, gin.H{
			"message":  "Account created successfully",
			"id":       user.ID,
			"name":     user.Name,
			"username": user.Username,
		})
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		req, ok := utils.BindJSON[LoginRequest](c)
		if !ok {
			return
		}

		req.Email = strings.ToLower(strings.TrimSpace(req.Email))

		if req.Email == "" || req.Password == "" {
			utils.BadRequest(c, "Email and password are required")
			return
		}

		user, err := gorm.G[models.User](db).
			Where("email = ?", req.Email).
			First(ctx)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.Unauthorized(c, "Invalid email or password")
			} else {
				utils.InternalError(c)
			}
			return
		}

		passwordMatch, err := argon2id.ComparePasswordAndHash(
			req.Password,
			user.PasswordHash,
		)
		if err != nil {
			utils.InternalError(c)
			return
		}

		if !passwordMatch {
			utils.Unauthorized(c, "Invalid email or password")
			return
		}

		// Generate a cryptographically secure token.
		tokenBytes := make([]byte, 32)
		if _, err := rand.Read(tokenBytes); err != nil {
			utils.InternalError(c)
			return
		}

		token := hex.EncodeToString(tokenBytes)

		// Store only the token hash.
		tokenHashBytes := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(tokenHashBytes[:])

		session := models.Session{
			UserID:    user.ID,
			TokenHash: tokenHash,
			ExpiresAt: time.Now().UTC().Add(sessionDuration),
		}

		if err := gorm.G[models.Session](db).Create(ctx, &session); err != nil {
			utils.InternalError(c)
			return
		}

		utils.ResponseOk(c, gin.H{
			"message":    "Login successful",
			"token":      token,
			"token_type": "Bearer",
			"expires_at": session.ExpiresAt,
		})
	}
}

func Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

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

		// Hash the token, since the database stores only its hash.
		hashBytes := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(hashBytes[:])

		// Delete the matching session to revoke the token.
		result, err := gorm.G[models.Session](db).
			Where("token_hash = ?", tokenHash).
			Delete(ctx)

		if err != nil {
			utils.InternalError(c)
			return
		}

		if result == 0 {
			utils.Unauthorized(c, "Invalid or expired session")
			return
		}

		utils.ResponseOk(c, gin.H{
			"message": "Logout successful",
		})
	}
}

func Me(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

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

		user, err := gorm.G[models.User](db).
			Where("id = ?", id).
			First(ctx)

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
