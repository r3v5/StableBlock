package models

import (
	"github.com/shopspring/decimal"
)

type Account struct {
	Address 		string 			 `gorm:"primaryKey;size:42" json:"address" binding:"required,len=42"`
	PasswordHash 	string 			 `gorm:"not null" json:"-"`
	IsValidator      bool   		 `gorm:"default:false" json:"is_validator"`
	IsZeroAddress    bool    		 `gorm:"default:false" json:"is_zero_address"`
	IsDepositAddress bool    		 `gorm:"default:false" json:"is_deposit_address"`
	SBBalance        decimal.Decimal `gorm:"type:numeric(20,8);default:0" json:"sb_balance"`
	TxSentCount      int             `gorm:"default:0" json:"tx_sent_count"`
}