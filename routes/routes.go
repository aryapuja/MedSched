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

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.JWTMiddleware())

	// Admin's routes
	protected.GET("/users", middleware.RoleRequired("admin"), handler.GetAllUsers)        // melihat semua user
	protected.GET("/users/by", middleware.RoleRequired("admin"), handler.GetUsersByQuery) // melihat user berdasarkan username/role
	protected.PUT("/users/by", middleware.RoleRequired("admin"), handler.UpdateUser)      // update data user lain
	protected.DELETE("/users/by", middleware.RoleRequired("admin"), handler.DeleteUser)   // hapus user lain

	// User's routes
	protected.GET("/users/me", middleware.RoleRequired("any"), handler.GetMyProfile)    // User bisa melihat profilnya sendiri
	protected.PUT("/users/me", middleware.RoleRequired("any"), handler.UpdateMyProfile) // User bisa mengupdate profilnya sendiri

}
