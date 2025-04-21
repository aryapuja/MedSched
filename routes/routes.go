package routes

import (
	"MedSched/handler"
	"MedSched/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/register", handler.RegisterUser)
	r.POST("/login", handler.LoginUser)
	r.GET("/doctor_schedule", handler.GetDoctorSchedule)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.JWTMiddleware())

	// Admin's routes
	protected.GET("/users", middleware.RoleRequired("admin"), handler.GetAllUsers)        // melihat semua user
	protected.GET("/users/by", middleware.RoleRequired("admin"), handler.GetUsersByQuery) // melihat user berdasarkan username/role
	protected.PUT("/users/by", middleware.RoleRequired("admin"), handler.UpdateUser)      // update data user lain
	protected.DELETE("/users/by", middleware.RoleRequired("admin"), handler.DeleteUser)   // hapus user lain

	// Doctor's routes
	protected.GET("/doctors", middleware.RoleRequired("admin"), handler.GetAllDoctors)      // Melihat daftar dokter
	protected.GET("/doctors/by", middleware.RoleRequired("admin"), handler.GetDoctorByID)   // Mencari dokter berdasarkan ID
	protected.POST("/doctors", middleware.RoleRequired("admin"), handler.CreateDoctor)      // Membuat dokter baru
	protected.PUT("/doctors/by", middleware.RoleRequired("admin"), handler.UpdateDoctor)    // Update dokter
	protected.DELETE("/doctors/by", middleware.RoleRequired("admin"), handler.DeleteDoctor) // Delete dokter

	// Schedule's routes
	protected.POST("/schedules", middleware.RoleRequired("admin"), handler.CreateSchedule)      // Membuat jadwal dokter baru
	protected.GET("/schedules", middleware.RoleRequired("admin"), handler.GetAllSchedules)      // Melihat semua jadwal dokter
	protected.GET("/schedules/by", middleware.RoleRequired("admin"), handler.GetScheduleByID)   // Mencari jadwal dokter berdasarkan ID
	protected.PUT("/schedules/by", middleware.RoleRequired("admin"), handler.UpdateSchedule)    // Update jadwal dokter
	protected.DELETE("/schedules/by", middleware.RoleRequired("admin"), handler.DeleteSchedule) // Delete jadwal dokter

	// Appoinment's routes
	protected.GET("/appointments", middleware.RoleRequired("admin"), handler.GetAllAppointments)             // Melihat semua appointment
	protected.GET("/appointments/date", middleware.RoleRequired("admin"), handler.GetAppointmentsWithDetail) // Melihat appointment berdasarkan tanggal

	// Schedule exceptions's routes
	protected.POST("/schedule_exceptions", middleware.RoleRequired("admin"), handler.CreateScheduleException)      // Menambah pengecualian jadwal
	protected.GET("/schedule_exceptions", middleware.RoleRequired("admin"), handler.GetAllScheduleExceptions)      // Melihat pengecualian jadwal
	protected.DELETE("/schedule_exceptions/by", middleware.RoleRequired("admin"), handler.DeleteScheduleException) // Hapus pengecualian jadwal

	// User's routes
	protected.GET("/users/me", middleware.RoleRequired("any"), handler.GetMyProfile)           // User bisa melihat profilnya sendiri
	protected.PUT("/users/me", middleware.RoleRequired("any"), handler.UpdateMyProfile)        // User bisa mengupdate profilnya sendiri
	protected.POST("/appointments", middleware.RoleRequired("any"), handler.CreateAppointment) // Membuat janji kontrol
}
