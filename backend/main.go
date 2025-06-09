package main

import (
	"log"
	"net/http"

	"backend/config"
	"backend/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database connection
	db, err := config.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize MinIO client
	minioClient, err := config.InitMinio()
	if err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	// Create router
	router := mux.NewRouter()

	// Setup routes
	routes.SetupRoutes(router)

	// Start the server
	port := ":8080"
	log.Printf("Backend service started on %s", port)
	log.Printf("Database connection established: %v", db != nil)
	log.Printf("MinIO client initialized: %v", minioClient != nil)
	log.Fatal(http.ListenAndServe(port, router))
}
