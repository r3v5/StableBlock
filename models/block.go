package models

import (
	"time"
)

type Block struct {
	Height       		int               `gorm:"primaryKey;autoIncrement;not null" json:"height"`
	Timestamp    		time.Time         `gorm:"autoCreateTime" json:"timestamp"`
	MaxTransactions     int               `gorm:"not null" json:"size"`
	ParentHash   		string            `json:"parent_hash"`
	Hash         		string            `gorm:"not null" json:"hash"`
}
