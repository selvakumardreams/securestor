package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/selvakumardreams/bluenoise/internal/utils"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
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

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Encrypt file content
	encryptedBytes, err := utils.Encrypt(fileBytes, utils.EncryptionKey)
	if err != nil {
		http.Error(w, "Failed to encrypt file", http.StatusInternalServerError)
		return
	}

	// Save encrypted file to the main bucket
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

	if _, err := dst.Write(encryptedBytes); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Asynchronously replicate the file to the replica buckets
	go utils.ReplicateFile(bucketName, header.Filename, encryptedBytes)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully: %s\n", header.Filename)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
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

	// Read encrypted file content
	encryptedBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Decrypt file content
	decryptedBytes, err := utils.Decrypt(encryptedBytes, utils.EncryptionKey)
	if err != nil {
		http.Error(w, "Failed to decrypt file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := w.Write(decryptedBytes); err != nil {
		http.Error(w, "Failed to download file", http.StatusInternalServerError)
		return
	}
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
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

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
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
