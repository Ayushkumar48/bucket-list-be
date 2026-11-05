package handlers

import (
	"bucket-list-be/src/models"
	"bucket-list-be/src/utils"
	"bucket-list-be/src/utils/schema"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OnboardingRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

func Onboarding(db *gorm.DB) gin.HandlerFunc {
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

		req, ok := utils.BindJSON[OnboardingRequest](c)
		if !ok {
			return
		}
		if !utils.ValidateZog(c, schema.OnboardingSchema.Validate(&req)) {
			return
		}

		name := req.Name
		username := req.Username
		bio := req.Bio
		result := db.Model(&models.User{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"name":     &name,
				"username": &username,
				"bio":      &bio,
			})
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
				utils.ConflictRequest(c, "Username already exists")
			} else {
				utils.InternalError(c)
			}
			return
		}
		if result.RowsAffected == 0 {
			utils.Unauthorized(c, "User not found")
			return
		}

		user, err := gorm.G[models.User](db).Where("id = ?", id).First(c.Request.Context())
		if err != nil {
			utils.InternalError(c)
			return
		}
		utils.ResponseOk(c, user)
	}
}
