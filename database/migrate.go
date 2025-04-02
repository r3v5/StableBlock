package database

import (
	"log"

	"github.com/r3v5/stableblock-api/models"
)

func Migrate() {
	err := DB.AutoMigrate(
		&models.Account{},
		&models.Stake{},
		&models.Block{},
		&models.Transaction{},
	)
	if err != nil {
		log.Fatal("❌ Migration failed:", err)
	}

	log.Println("✅ All tables migrated successfully.")
}
