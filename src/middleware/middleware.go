package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"bucket-list-be/src/models"
	"bucket-list-be/src/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthRequired(db *gorm.DB) gin.HandlerFunc {
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

		// Hash the raw token exactly as we did during login.
		hashBytes := sha256.Sum256([]byte(token))
		tokenHash := hex.EncodeToString(hashBytes[:])

		// Only accept sessions that have not expired.
		session, err := gorm.G[models.Session](db).
			Where("token_hash = ? AND expires_at > ?", tokenHash, time.Now().UTC()).
			First(c.Request.Context())

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.Unauthorized(c, "Invalid or expired session")
			} else {
				utils.InternalError(c)
			}
			return
		}

		// Make the authenticated user ID available to downstream handlers.
		c.Set("userID", session.UserID)

		c.Next()
	}
}
