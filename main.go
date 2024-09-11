package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const storageDir = "./storage"

// List of buckets for replication
var replicationBuckets = []string{"replica1", "replica2"}

func main() {
	http.HandleFunc("/create-bucket", createBucketHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/delete", deleteHandler) // DELETE method

	fmt.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}

func createBucketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.URL.Query().Get("bucket")
	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	bucketPath := filepath.Join(storageDir, bucketName)
	if err := os.MkdirAll(bucketPath, os.ModePerm); err != nil {
		http.Error(w, "Failed to create bucket", http.StatusInternalServerError)
		return
	}

	// Replicate the bucket
	for _, replica := range replicationBuckets {
		replicaPath := filepath.Join(storageDir, replica, bucketName)
		if err := os.MkdirAll(replicaPath, os.ModePerm); err != nil {
			http.Error(w, "Failed to create replica bucket", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Bucket created successfully: %s\n", bucketName)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.URL.Query().Get("bucket")
	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save file to the main bucket
	bucketPath := filepath.Join(storageDir, bucketName)
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		http.Error(w, "Bucket does not exist", http.StatusBadRequest)
		return
	}

	dst, err := os.Create(filepath.Join(bucketPath, header.Filename))
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Asynchronously replicate the file to the replica buckets
	go replicateFile(bucketName, header.Filename, file)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully: %s\n", header.Filename)
}

func replicateFile(bucketName, filename string, file io.Reader) {
	for _, replica := range replicationBuckets {
		replicaPath := filepath.Join(storageDir, replica, bucketName)
		if _, err := os.Stat(replicaPath); os.IsNotExist(err) {
			fmt.Printf("Replica bucket %s does not exist\n", replica)
			continue
		}

		replicaDst, err := os.Create(filepath.Join(replicaPath, filename))
		if err != nil {
			fmt.Printf("Failed to create file in replica bucket %s: %v\n", replica, err)
			continue
		}
		defer replicaDst.Close()

		if _, err := io.Copy(replicaDst, file); err != nil {
			fmt.Printf("Failed to save file in replica bucket %s: %v\n", replica, err)
		}
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.URL.Query().Get("bucket")
	filename := r.URL.Query().Get("filename")
	if bucketName == "" || filename == "" {
		http.Error(w, "Bucket name and filename are required", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(storageDir, bucketName, filename)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
		return
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.URL.Query().Get("bucket")
	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	bucketPath := filepath.Join(storageDir, bucketName)
	files, err := os.ReadDir(bucketPath)
	if err != nil {
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		fmt.Fprintln(w, file.Name())
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	bucketName := r.URL.Query().Get("bucket")
	filename := r.URL.Query().Get("filename")
	if bucketName == "" || filename == "" {
		http.Error(w, "Bucket name and filename are required", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(storageDir, bucketName, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	if err := os.Remove(filePath); err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	// Delete the file from the replica buckets
	for _, replica := range replicationBuckets {
		replicaFilePath := filepath.Join(storageDir, replica, bucketName, filename)
		if _, err := os.Stat(replicaFilePath); !os.IsNotExist(err) {
			if err := os.Remove(replicaFilePath); err != nil {
				http.Error(w, "Failed to delete file from replica bucket", http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File deleted successfully: %s\n", filename)
}
