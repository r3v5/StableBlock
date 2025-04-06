package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/r3v5/stableblock-api/database"
	"github.com/r3v5/stableblock-api/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func HandleGetSBTokens(c *gin.Context) {
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

	// Transactionally credit 3 SB tokens
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		var account models.Account

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&account, "address = ?", address).Error; err != nil {
			c.JSON(404, gin.H{"error": "Account not found"})
			return err
		}

		account.SBBalance = account.SBBalance.Add(decimal.NewFromInt(3))

		if err := tx.Save(&account).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update balance"})
			return err
		}

		c.JSON(200, gin.H{
			"address":    account.Address,
			"new_balance": account.SBBalance,
			"message":    "3 SB tokens received successfully",
		})

		return nil
	})

	if err != nil {
		c.JSON(500, gin.H{"error": "Transaction failed"})
	}
}
