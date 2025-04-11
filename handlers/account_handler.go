package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/r3v5/stableblock-api/database"
	"github.com/r3v5/stableblock-api/models"
)

func HandleGetAccount(c *gin.Context) {
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

	account := &models.Account{}
	if err := database.DB.First(account, "address = ?", address).Error; err != nil {
		c.JSON(404, gin.H{"error": "Account not found"})
		return
	}

	c.JSON(200, gin.H{
		"address":            account.Address,
		"name":				  account.Name,	
		"sb_balance":         account.SBBalance,
		"tx_sent_count":      account.TxSentCount,
		"date_created":		  account.DateCreated,	
	})
}