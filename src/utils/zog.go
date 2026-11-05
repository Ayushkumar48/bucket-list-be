package utils

import (
	"strings"

	z "github.com/Oudwins/zog"
	"github.com/gin-gonic/gin"
)

func ValidateZog(c *gin.Context, issues z.ZogIssueList) bool {
	if len(issues) == 0 {
		return true
	}

	validationErrors := make([]gin.H, 0, len(issues))

	for _, issue := range issues {
		validationErrors = append(validationErrors, gin.H{
			"field":   strings.Join(issue.Path, "."),
			"message": issue.Message,
		})
	}

	ValidationError(c, validationErrors)
	return false
}
