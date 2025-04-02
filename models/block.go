package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Block struct {
	ID           		uint              `gorm:"primaryKey" json:"id"`
	Height       		int               `gorm:"uniqueIndex;autoIncrement;default:0" json:"height"`
	Timestamp    		time.Time         `gorm:"autoCreateTime" json:"timestamp"`
	Size         		int               `gorm:"not null" json:"size"`
	ParentHash   		string            `json:"parent_hash"`
	Hash         		string            `gorm:"not null" json:"hash"`
	FeeRecipientAddress string     		  `gorm:"type:char(42);not null" json:"fee_recipient"`
	FeeRecipient 		Account           `gorm:"foreignKey:FeeRecipientAddress;references:Address" json:"-"`
	BlockReward 		decimal.Decimal   `gorm:"type:numeric(20,8);not null" json:"block_reward"`
}
