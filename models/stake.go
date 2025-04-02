package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Stake struct {
	ID             uint      		`gorm:"primaryKey" json:"id"`
	AccountAddress string    		`gorm:"type:char(42);not null" json:"account_address"`
	Account 	   Account			`gorm:"foreignKey:AccountAddress;references:Address" json:"-"`
	Amount         decimal.Decimal	`gorm:"type:numeric(20,8);not null" json:"amount"`
	Timestamp      time.Time 		`gorm:"autoCreateTime" json:"timestamp"`
}
