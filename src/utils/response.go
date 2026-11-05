package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Unauthorized(c *gin.Context, msg ...string) {
	m := "unauthorized"
	if len(msg) > 0 {
		m = msg[0]
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": m})
}

func InternalError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
}

func BadRequest(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
}

func ForbiddenRequest(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
}
func ConflictRequest(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusConflict, gin.H{
		"error": msg,
	})
}

func ResponseOk(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func ValidationError(c *gin.Context, issues any) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error":  "Validation failed",
		"issues": issues,
	})
}
