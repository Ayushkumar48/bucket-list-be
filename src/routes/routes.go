package routes

import (
	"bucket-list-be/src/config"
	"bucket-list-be/src/handlers"
	"bucket-list-be/src/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Setup(db *gorm.DB, cfg config.Config) *gin.Engine {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	auth := router.Group("/auth")
	auth.POST("/login", handlers.Login(db))
	auth.POST("/signup", handlers.Signup(db))
	auth.POST("/logout", handlers.Logout(db))

	protected := router.Group("/")
	protected.Use(middleware.AuthRequired(db))
	{
		protected.GET("/me", handlers.Me(db))
		// protected.POST("/bucket-items", handlers.CreateBucketItem)
		// protected.GET("/bucket-items", handlers.ListBucketItems)
	}

	return router
}
