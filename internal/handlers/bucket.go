package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

const storageDir = "../../storage"

// List of buckets for replication
var replicationBuckets = []string{"replica1", "replica2"}

func CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
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
