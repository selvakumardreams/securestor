package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

const storageDir = "./storage"

// List of buckets for replication
var replicationBuckets = []string{"replica1", "replica2"}

func ReplicateFile(bucketName, filename string, fileBytes []byte) {
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
