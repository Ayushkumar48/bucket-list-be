package routes

import (
	"bucket-list-be/src/config"
	"bucket-list-be/src/handlers"
	"bucket-list-be/src/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, redisClient *redis.Client, cfg config.Config) *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	auth := router.Group("/auth")
	auth.POST("/login", handlers.RequestMagicLink(db, redisClient, cfg))
	auth.POST("/verify", handlers.VerifyMagicLink(db, redisClient))

	protected := router.Group("/")
	protected.Use(middleware.AuthRequired(db))
	{
		protected.POST("/auth/logout", handlers.Logout(db))
		protected.GET("/me", handlers.Me(db))
		protected.PATCH("/onboarding", handlers.Onboarding(db))
		// protected.POST("/bucket-items", handlers.CreateBucketItem)
		// protected.GET("/bucket-items", handlers.ListBucketItems)
	}

	return router
}
