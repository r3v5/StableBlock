package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/r3v5/stableblock-api/database"
	"github.com/r3v5/stableblock-api/models"
)


func HandleGetBlocks(c *gin.Context) {
	var blocks []models.Block

	if err := database.DB.Order("height ASC").Find(&blocks).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch blocks"})
		return
	}

	c.JSON(200, gin.H{"blocks": blocks})
}


func HandleGetBlock(c *gin.Context) {
	heightParam := c.Param("height")

	var block models.Block
	if err := database.DB.First(&block, "height = ?", heightParam).Error; err != nil {
		c.JSON(404, gin.H{"error": "Block not found"})
		return
	}

	c.JSON(200, gin.H{"block": block})
}
