package main

import (
	"fmt"
	"net/http"

	"github.com/selvakumardreams/bluenoise/internal/handlers"
)

func main() {
	http.HandleFunc("/create-bucket", handlers.CreateBucketHandler) // create a new bucket
	http.HandleFunc("/upload", handlers.UploadHandler)              // upload a file to a bucket
	http.HandleFunc("/download", handlers.DownloadHandler)          // download a file from a bucket
	http.HandleFunc("/list", handlers.ListHandler)                  // List files in a bucket
	http.HandleFunc("/delete", handlers.DeleteHandler)              // Delete a file from a bucket
	http.HandleFunc("/health", handlers.HealthCheckHandler)         // Health check

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
