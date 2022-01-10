package util

import (
	"os"
)

const (
	StorageApiSecretId  = "STORAGE_API_SECRET_ID"
	StorageApiSecretKey = "STORAGE_API_SECRET_KEY"
)

func GetSercet() (secretID, secretKey string) {
	secretID = os.Getenv(StorageApiSecretId)
	secretKey = os.Getenv(StorageApiSecretKey)
	return
}
