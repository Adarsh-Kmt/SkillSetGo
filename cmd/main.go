package main

import (
	"log"
	"net/http"

	db "github.com/adarsh-kmt/skillsetgo/pkg/db/config"
	"github.com/adarsh-kmt/skillsetgo/pkg/handler"
	"github.com/adarsh-kmt/skillsetgo/pkg/middleware"
	"github.com/adarsh-kmt/skillsetgo/pkg/service"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	if err := db.PostgresDBClientInit(); err != nil {
		log.Printf("Database initialization error: %v", err)
		log.Fatal(err)
	}

	// Initialize services
	jobService := service.NewJobServiceImpl()
	studentService := service.NewStudentServiceImpl()
	authService := service.NewAuthServiceImpl()

	// Initialize handlers
	jobHandler := handler.NewJobHandler(jobService)
	studentHandler := handler.NewStudentHandler(studentService)
	authHandler := handler.NewAuthHandler(authService)

	// Create router
	router := mux.NewRouter()

	// Add CORS middleware
	router.Use(middleware.CorsMiddleware)

	// Setup routes
	router = jobHandler.MuxSetup(router)
	router = studentHandler.MuxSetup(router)
	router = authHandler.MuxSetup(router)

	// Start server
	log.Printf("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
