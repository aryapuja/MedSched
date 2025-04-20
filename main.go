package main

import (
	"fmt"
	"log"
	"os"

	"MedSched/config"
	_ "MedSched/docs"
	"MedSched/routes"
	"MedSched/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {

	//file .env
	if os.Getenv("RAILWAY_ENVIRONMENT_NAME") == "" { 
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	// DB
	db, err := config.SetupDatabase()
	if err != nil {
		log.Fatal("Error setting up database: ", err)
	}
	defer db.Close()

	// Router Gin
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Setup Routes
	routes.SetupRoutes(r)

	// scheduler pengingat
	services.StartReminderScheduler()
	// log.Println("Scheduler started")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Jalankan aplikasi
	err = r.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
