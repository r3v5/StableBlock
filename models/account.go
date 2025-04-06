package models

import (
	"github.com/shopspring/decimal"
)

type Account struct {
	Address 		string 			 `gorm:"primaryKey;size:42" json:"address" binding:"required,len=42"`
	Name 			string 			 `gorm:"type:varchar(32);not null" json:"name" binding:"required"`
	PasswordHash 	string 			 `gorm:"not null" json:"-"`
	SBBalance       decimal.Decimal  `gorm:"type:numeric(20,8);default:0" json:"sb_balance"`
	TxSentCount     int              `gorm:"default:0" json:"tx_sent_count"`
	RefreshToken 	*string 		 `gorm:"type:text"`
}