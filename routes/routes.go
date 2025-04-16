package routes

import (
	"MedSched/handler"
	"MedSched/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes untuk mendefinisikan semua rute API
func SetupRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/register", handler.RegisterUser)
	r.POST("/login", handler.LoginUser)

	// Protected routes (requires JWT)
	protected := r.Group("/api")
	protected.Use(middleware.JWTMiddleware())
	{
		// Tambahkan route yang membutuhkan autentikasi JWT
	}
}
