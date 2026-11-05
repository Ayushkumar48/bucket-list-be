package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindJSON[T any](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return req, false
	}
	return req, true
}
