package main

import (
	"fmt"
	"net/http"

	"github.com/selvakumardreams/bluenoise/internal/handlers"
)

func main() {
	http.HandleFunc("/create-bucket", corsMiddleware(handlers.CreateBucketHandler))           // create a new bucket
	http.HandleFunc("/upload", corsMiddleware(handlers.UploadHandler))                        // upload a file to a bucket
	http.HandleFunc("/download", corsMiddleware(handlers.DownloadHandler))                    // download a file from a bucket
	http.HandleFunc("/list", corsMiddleware(handlers.ListHandler))                            // List files in a bucket
	http.HandleFunc("/delete", corsMiddleware(handlers.DeleteHandler))                        // Delete a file from a bucket
	http.HandleFunc("/health", corsMiddleware(handlers.HealthCheckHandler))                   // Health check
	http.HandleFunc("/search", corsMiddleware(handlers.SearchHandler))                        // Search for a file in a bucket
	http.HandleFunc("/update-metadata", corsMiddleware(handlers.UpdateCustomMetadataHandler)) // Update custom metadata

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", corsMiddleware(fs.ServeHTTP))

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}
