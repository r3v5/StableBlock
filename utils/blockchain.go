package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/r3v5/stableblock-api/database"
	"github.com/r3v5/stableblock-api/models"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"
)


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

func Keccak256Hash(data string) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(data))
	return "0x" + hex.EncodeToString(hash.Sum(nil))
}

func GenerateTransactionHash(from, to string, value decimal.Decimal, timestamp time.Time) string {
	data := fmt.Sprintf("%s:%s:%s:%d", from, to, value.String(), timestamp.UnixNano())
	return Keccak256Hash(data)
}

func GenerateBlockHash(height int, parentHash string, timestamp time.Time) string {
	data := fmt.Sprintf("%d:%s:%d", height, parentHash, timestamp.UnixNano())
	return Keccak256Hash(data)
}


func GetOrCreateBlockWithFreeSlot(tx *gorm.DB, defaultMaxTx int) (*models.Block, error) {
	var blocks []models.Block

	// 1. Load all blocks ordered by height
	if err := tx.Order("height ASC").Find(&blocks).Error; err != nil {
		return nil, err
	}

	// 2. Loop over them to find one with available space
	for _, block := range blocks {
		var count int64
		if err := tx.Model(&models.Transaction{}).
			Where("block_height = ?", block.Height).
			Count(&count).Error; err != nil {
			return nil, err
		}

		if int(count) < block.MaxTransactions {
			return &block, nil
		}
	}

	// 3. All blocks are full — create a new block now
	var latestBlock models.Block
	if err := tx.Order("height DESC").First(&latestBlock).Error; err != nil {
		// fallback to genesis
		latestBlock.Height = -1
		latestBlock.Hash = "0x0000000000000000000000000000000000000000000000000000000000000000"
	}

	newBlock := models.Block{
		ParentHash:      latestBlock.Hash,
		MaxTransactions: defaultMaxTx,
		Timestamp:       time.Now(),
	}
	newBlock.Hash = GenerateBlockHash(latestBlock.Height+1, newBlock.ParentHash, newBlock.Timestamp)

	if err := tx.Create(&newBlock).Error; err != nil {
		return nil, err
	}

	return &newBlock, nil
}

