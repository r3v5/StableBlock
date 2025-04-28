package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/r3v5/stableblock-api/models"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
	"gorm.io/gorm"
)

type BlockchainUtil interface {
	GenerateUniqueAddress() (string, error)
	Keccak256Hash(data string) string
	GenerateTransactionHash(from, to string, value decimal.Decimal, timestamp time.Time) string
	GenerateBlockHash(height int, parentHash string, timestamp time.Time) string
	GetOrCreateBlockWithFreeSlot(tx *gorm.DB, defaultMaxTx int) (*models.Block, error)
}

type DefaultBlockchainUtil struct {
	DB *gorm.DB
}

func (u *DefaultBlockchainUtil) GenerateUniqueAddress() (string, error) {
	const maxAttempts = 10

	for i := 0; i < maxAttempts; i++ {
		bytes := make([]byte, 20)
		_, err := rand.Read(bytes)
		if err != nil {
			return "", err
		}
		address := "0x" + hex.EncodeToString(bytes)

		var existing models.Account
		result := u.DB.First(&existing, "address = ?", address)
		if result.Error == nil {
			continue
		}
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return address, nil
		}
		return "", fmt.Errorf("DB error: %v", result.Error)
	}

	return "", errors.New("could not generate unique address after multiple attempts")
}

func (u *DefaultBlockchainUtil) Keccak256Hash(data string) string {
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(data))
	return "0x" + hex.EncodeToString(hash.Sum(nil))
}

func (u *DefaultBlockchainUtil) GenerateTransactionHash(from, to string, value decimal.Decimal, timestamp time.Time) string {
	data := fmt.Sprintf("%s:%s:%s:%d", from, to, value.String(), timestamp.UnixNano())
	return u.Keccak256Hash(data)
}

func (u *DefaultBlockchainUtil) GenerateBlockHash(height int, parentHash string, timestamp time.Time) string {
	data := fmt.Sprintf("%d:%s:%d", height, parentHash, timestamp.UnixNano())
	return u.Keccak256Hash(data)
}


func (u *DefaultBlockchainUtil) CreateNewBlock(tx *gorm.DB, defaultMaxTx int) (*models.Block, error) {
	var latestBlock models.Block
	err := tx.Order("height DESC").First(&latestBlock).Error

	if err != nil {
		return nil, err
	}

	newBlock := models.Block{
		ParentHash:      latestBlock.Hash,
		MaxTransactions: defaultMaxTx,
		Timestamp:       time.Now(),
	}
	newBlock.Hash = u.GenerateBlockHash(latestBlock.Height+1, newBlock.ParentHash, newBlock.Timestamp)

	if err := tx.Create(&newBlock).Error; err != nil {
		return nil, err
	}

	return &newBlock, nil
}

func (u *DefaultBlockchainUtil) GetBlockWithFreeSlot(tx *gorm.DB) (*models.Block, error) {
	var blocks []models.Block
	if err := tx.Order("height ASC").Find(&blocks).Error; err != nil {
		return nil, err
	}

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

	return nil, nil
}