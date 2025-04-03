package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/r3v5/stableblock-api/database"
	"github.com/r3v5/stableblock-api/models"
	"github.com/r3v5/stableblock-api/services"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Address  string `json:"address" binding:"required,len=42"`
	Password string `json:"password" binding:"required"`
}


func HandlePostRegister(c *gin.Context) {
	inputPassword := &RegisterInput{}
	inputError := c.ShouldBindJSON(inputPassword)
	if inputError != nil {
		c.JSON(400, gin.H{"error": inputError.Error()})
		return
	}

	address, err := services.GenerateUniqueAddress()
	if err != nil {
		c.JSON(400, gin.H{"error": "Could not generate address, try again"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(inputPassword.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(400, gin.H{"error": "Could not hash password"})
		return
	}

	account := &models.Account{
		Address:        address,
		PasswordHash:   string(hashedPassword),
		IsValidator:    false,
		IsZeroAddress:  false,
		IsDepositAddress: false,
		SBBalance:      decimal.NewFromFloat(10.0),
		TxSentCount:    0,
	}

	if err := database.DB.Create(account).Error; err != nil {
		c.JSON(400, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(201, gin.H{
		"address": address,
	})
}


func HandlePostLogin(c *gin.Context) {
	loginData := &LoginInput{}
	loginDataError := c.ShouldBindJSON(loginData)

	if loginDataError != nil {
		c.JSON(400, gin.H{"error": loginDataError.Error()})
		return
	}

	account := &models.Account{}
	if err := database.DB.First(account, "address = ?", loginData.Address).Error; err != nil {
		c.JSON(401, gin.H{"error": "Address is not found in StableBlock"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(loginData.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Incorrect password"})
		return
	}

	accessToken, refreshToken, err := services.GenerateTokens(account.Address)
	if err != nil {
		c.JSON(400, gin.H{"error": "Could not generate tokens"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}