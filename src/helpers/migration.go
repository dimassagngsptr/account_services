package helpers

import (
	"account_services/src/configs"
	"account_services/src/models"
	"log"
)

func Migration() {
	log.Println("Running Migrations")
	err := configs.DB.AutoMigrate(
		&models.User{},
		&models.Account{},
	)
	if err != nil {
		log.Fatalf("can't running migration")
	}
}
