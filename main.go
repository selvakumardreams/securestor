package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const storageDir = "./storage"
const encryptionKey = "Yf90Frf3DorOqeDfK4VGRIeQfGKUgkle" // 32 bytes key for AES-256

// List of buckets for replication
var replicationBuckets = []string{"replica1", "replica2"}

func main() {
	http.HandleFunc("/create-bucket", createBucketHandler) // create a new bucket
	http.HandleFunc("/upload", uploadHandler)              // upload a file to a bucket
	http.HandleFunc("/download", downloadHandler)          // download a file from a bucket
	http.HandleFunc("/list", listHandler)                  // List files in a bucket
	http.HandleFunc("/delete", deleteHandler)              // Delete a file from a bucket

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

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Encrypt file content
	encryptedBytes, err := encrypt(fileBytes, encryptionKey)
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
	go replicateFile(bucketName, header.Filename, encryptedBytes)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully: %s\n", header.Filename)
}

func replicateFile(bucketName, filename string, fileBytes []byte) {
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

		if _, err := replicaDst.Write(fileBytes); err != nil {
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

	// Read encrypted file content
	encryptedBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Decrypt file content
	decryptedBytes, err := decrypt(encryptedBytes, encryptionKey)
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

func encrypt(data []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(passphrase))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func decrypt(data []byte, passphrase string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(passphrase))
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
