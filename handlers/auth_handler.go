package handlers

import (
	"fmt"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/r3v5/stableblock-api/database"
	"github.com/r3v5/stableblock-api/models"
	"github.com/r3v5/stableblock-api/utils"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Password string `json:"password" binding:"required,min=6"`
	Name 	 string	`json:"name" binding:"required"`
}

type LogineRequest struct {
	Address  string `json:"address" binding:"required,len=42"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}


func HandlePostRegister(c *gin.Context) {
	registerData := &RegisterRequest{}
	inputError := c.ShouldBindJSON(registerData)
	util := &utils.DefaultBlockchainUtil{DB: database.DB}
	if inputError != nil {
		c.JSON(400, gin.H{"error": inputError.Error()})
		return
	}

	address, err := util.GenerateUniqueAddress()
	if err != nil {
		c.JSON(400, gin.H{"error": "Could not generate address, try again"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(400, gin.H{"error": "Could not hash password"})
		return
	}

	account := &models.Account{
		Address:        address,
		Name:			registerData.Name,
		PasswordHash:   string(hashedPassword),
		SBBalance:      decimal.NewFromFloat(10.0),
		TxSentCount:    0,
		RefreshToken:   nil,
		DateCreated:    time.Now(),
	}

	if err := database.DB.Create(account).Error; err != nil {
		c.JSON(400, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(201, gin.H{
		"address": address,
		"name": registerData.Name,
	})
}


func HandlePostLogin(c *gin.Context) {
	loginData := &LogineRequest{}
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

	accessToken, refreshToken, err := utils.GenerateTokens(account.Address)
	if err != nil {
		c.JSON(400, gin.H{"error": "Could not generate tokens"})
		return
	}

	account.RefreshToken = &refreshToken
	if err := database.DB.Save(account).Error; err != nil {
		c.JSON(400, gin.H{"error": "Could not store refresh token"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}


func HandlePostRefresh(c *gin.Context) {
	input := &RefreshRequest{}
	if err := c.ShouldBindJSON(input); err != nil {
		c.JSON(400, gin.H{"error": "Refresh token required"})
		return
	}

	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return utils.GetSecret(), nil
	})

	var address string
	if token != nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if addr, ok := claims["address"].(string); ok {
				address = addr
			}
		}
	}

	if err != nil || !token.Valid {
		if address != "" {
			account := &models.Account{}
			if err := database.DB.First(account, "address = ?", address).Error; err == nil {
				account.RefreshToken = nil
				if err := database.DB.Save(account).Error; err != nil {
					log.Printf("Failed to update refresh token for address %s: %v", address, err)
				}
				
			}
		}

		c.JSON(401, gin.H{"error": "Refresh token is no longer valid"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(401, gin.H{"error": "Invalid token claims"})
		return
	}

	address, ok = claims["address"].(string)
	if !ok || address == "" {
		c.JSON(401, gin.H{"error": "Invalid or missing address in token"})
		return
	}

	account := &models.Account{}
	if err := database.DB.First(account, "address = ?", address).Error; err != nil {
		c.JSON(401, gin.H{"error": "Account not found"})
		return
	}

	if account.RefreshToken == nil || *account.RefreshToken != input.RefreshToken {
		account.RefreshToken = nil
		if err := database.DB.Save(account).Error; err != nil {
			log.Printf("Failed to update refresh token for address %s: %v", address, err)
		}
		c.JSON(401, gin.H{"error": "Refresh token is no longer valid"})
		return
	}

	newAccessToken, err := utils.GenerateAccessToken(address)
	if err != nil {
		c.JSON(400, gin.H{"error": "Could not generate new access token"})
		return
	}

	c.JSON(200, gin.H{
		"access_token": newAccessToken,
	})
}


func HandlerPostLogout(c *gin.Context) {
	address := c.GetString("address")

	account := &models.Account{}

	if err := database.DB.First(account, "address = ?", address).Error; err != nil {
		c.JSON(401, gin.H{"error": "Account not found"})
		return
	}

	account.RefreshToken = nil

	if err := database.DB.Save(account).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to log out"})
		return
	}

	c.JSON(200, gin.H{"message": "Logged out successfully"})
}
