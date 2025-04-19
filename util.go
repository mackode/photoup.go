package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"
)

func uploadDir() (string, error) {
	root := os.Getenv("DOCUMENT_ROOT")
	if root == "" {
		return "", fmt.Errorf("DOCUMENT_ROOT not set")
	}

	randomStr := make([]byte, 32)
	if _, err := rand.Read(randomStr); err != nil {
		return "", err
	}
	hash := sha256.Sum256(randomStr)
	shaDir := hex.EncodeToString(hash[:])

	dir := path.Join(root, "uploads")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
		err := os.WriteFile(path.Join(dir, ".htaccess"), []byte("Option-Indexes\n"), 0644)
		if err != nil {
			return "", err
		}
	}
	path := path.Join(dir, shaDir)
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return "", err
	}
	return path, nil
}
