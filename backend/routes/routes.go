package routes

import (
	"backend/controllers"
	"net/http"

	"github.com/gorilla/mux"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *mux.Router) {
	// API Routes
	api := r.PathPrefix("/api").Subrouter()
	
	// File routes
	fileRoutes := api.PathPrefix("/files").Subrouter()
	fileRoutes.HandleFunc("", controllers.GetAllFiles).Methods(http.MethodGet)
	fileRoutes.HandleFunc("/{id:[0-9]+}", controllers.GetFileByID).Methods(http.MethodGet)
	fileRoutes.HandleFunc("", controllers.UploadFile).Methods(http.MethodPost) // Keeping single-part for small files
	fileRoutes.HandleFunc("/{id:[0-9]+}", controllers.DeleteFile).Methods(http.MethodDelete)
	fileRoutes.HandleFunc("/{id:[0-9]+}/download", controllers.DownloadFile).Methods(http.MethodGet)
	
	// Multi-part upload routes
	mpRoutes := api.PathPrefix("/uploads").Subrouter()
	mpRoutes.HandleFunc("/initiate", controllers.InitiateMultipartUpload).Methods(http.MethodPost)
	mpRoutes.HandleFunc("/{uploadId}/parts", controllers.UploadPart).Methods(http.MethodPost)
	mpRoutes.HandleFunc("/{uploadId}/complete", controllers.CompleteMultipartUpload).Methods(http.MethodPost)
	mpRoutes.HandleFunc("/{uploadId}/abort", controllers.AbortMultipartUpload).Methods(http.MethodDelete)
	
	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	}).Methods(http.MethodGet)
	
	// Global middleware for CORS
	r.Use(corsMiddleware)
}

// corsMiddleware handles CORS configuration
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}
