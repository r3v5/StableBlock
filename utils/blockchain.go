package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/r3v5/stableblock-api/database"
	"github.com/r3v5/stableblock-api/models"
	"gorm.io/gorm"
)

type Blockchain interface {
	GenerateUniqueAddress() (string, error)
}


func GenerateUniqueAddress() (string, error) {
	const maxAttempts = 10

	for i := 0; i < maxAttempts; i++ {
		// Generate a random 20-byte address
		bytes := make([]byte, 20)
		_, err := rand.Read(bytes)
		if err != nil {
			return "", err
		}
		address := "0x" + hex.EncodeToString(bytes)

		// Check if address already exists in the DB
		var existing models.Account
		result := database.DB.First(&existing, "address = ?", address)
		if result.Error == nil {
			continue
		}

		if result.Error.Error() == "record not found" || errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return address, nil
		}

		return "", fmt.Errorf("DB error: %v", result.Error)
	}

	return "", errors.New("could not generate unique address after multiple attempts")
}