package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	TransactionHash string          `gorm:"primaryKey;type:char(66)" json:"transaction_hash"`
	BlockHeight     int             `gorm:"not null" json:"block_height"`
	Block           Block           `gorm:"foreignKey:BlockHeight;references:Height" json:"-"`
	Timestamp       time.Time       `gorm:"autoCreateTime" json:"timestamp"`
	FromAddress     string          `gorm:"type:char(42);not null" json:"from_address"`
	FromAccount     Account         `gorm:"foreignKey:FromAddress;references:Address" json:"-"`
	ToAddress       string          `gorm:"type:char(42);not null" json:"to_address"`
	ToAccount       Account         `gorm:"foreignKey:ToAddress;references:Address" json:"-"`
	Value           decimal.Decimal `gorm:"type:numeric(20,8);not null" json:"value"`
}
