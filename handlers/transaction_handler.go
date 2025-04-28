package handlers

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/r3v5/stableblock-api/database"
	"github.com/r3v5/stableblock-api/models"
	"github.com/r3v5/stableblock-api/utils"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)


type TransactionRequest struct {
	ToAddress string          `json:"to_address" binding:"required"`
	Value     decimal.Decimal `json:"value" binding:"required"` 
}

func HandlePostTransaction(c *gin.Context) {
	addr, exists := c.Get("address")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	address, ok := addr.(string)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	incomingTxData := &TransactionRequest{}
	if err := c.ShouldBindJSON(incomingTxData); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if incomingTxData.Value.LessThan(decimal.NewFromFloat(0.5)) {
		c.JSON(400, gin.H{"error": "value must be at least 0.5"})
		return
	}

	// Transactional logic
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var sender models.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&sender, "address = ?", address).Error; err != nil {
			return err
		}

		if sender.SBBalance.LessThan(incomingTxData.Value) {
			c.JSON(400, gin.H{"error": "Insufficient balance"})
			return errors.New("abort")
		}

		var recipient models.Account
		if err := tx.First(&recipient, "address = ?", incomingTxData.ToAddress).Error; err != nil {
			c.JSON(400, gin.H{"error": "Recipient not found"})
			return errors.New("abort")
		}

		// Prepare block
		util := &utils.DefaultBlockchainUtil{DB: database.DB}

		block, err := util.GetBlockWithFreeSlot(tx)
		if err != nil {
			c.JSON(400, gin.H{"error": "Unable to allocate block"})
			return err
		}

		if block == nil {
			block, err = util.CreateNewBlock(tx, 3)
			if err != nil {
				c.JSON(400, gin.H{"error": "Unable to create new block"})
				return err
			}
		}

		now := time.Now()
		txHash := util.GenerateTransactionHash(sender.Address, recipient.Address, incomingTxData.Value, now)

		transaction := models.Transaction{
			TransactionHash: txHash,
			BlockHeight:     block.Height,
			Timestamp:       now,
			FromAddress:     sender.Address,
			ToAddress:       recipient.Address,
			Value:           incomingTxData.Value,
		}

		sender.SBBalance = sender.SBBalance.Sub(incomingTxData.Value)
		sender.TxSentCount++
		recipient.SBBalance = recipient.SBBalance.Add(incomingTxData.Value)

		if err := tx.Save(&sender).Error; err != nil {
			return err
		}
		if err := tx.Save(&recipient).Error; err != nil {
			return err
		}
		if err := tx.Create(&transaction).Error; err != nil {
			return err
		}

		c.JSON(200, gin.H{
			"transaction_hash": txHash,
			"block_height":     block.Height,
			"timestamp":        now,
			"from_address":		sender.Address,
			"to_address": 		recipient.Address,
			"value": 			incomingTxData.Value,
		})
		return nil
	})

	if err != nil && err.Error() != "abort" {
		c.JSON(400, gin.H{"error": "Transaction failed"})
	}
}



func HandleGetTransactions(c *gin.Context) {
	var transactions []models.Transaction

	if err := database.DB.Order("timestamp ASC").Find(&transactions).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	c.JSON(200, gin.H{"transactions": transactions})
}


func HandleGetTransaction(c *gin.Context) {
	hash := c.Param("hash")

	var tx models.Transaction
	if err := database.DB.First(&tx, "transaction_hash = ?", hash).Error; err != nil {
		c.JSON(404, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(200, gin.H{"transaction": tx})
}
