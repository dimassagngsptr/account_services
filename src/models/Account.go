package models

import (
	"account_services/src/configs"

	"github.com/google/uuid"
)

type Account struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	AccountNumber string    `json:"account_number"`
	Nominal       float64   `json:"nominal"`
	UserID        string    `json:"user_id"`
	User          User      `gorm:"foreignKey:UserID"`
}

type AccountRequest struct {
	AccountNumber string  `json:"account_number"`
	Nominal       float64 `json:"nominal"`
}

func CreateAccount(newAccount *Account) (*Account, error) {
	if err := configs.DB.Create(&newAccount).Error; err != nil {
		return nil, err
	}
	return newAccount, nil
}

func FindByAccNumber(accNumber string) *Account {
	var account Account
	result := configs.DB.Where("account_number = ?", accNumber).First(&account)
	if result.RowsAffected == 0 {
		return nil
	}
	return &account
}

func UpdateBalance(accNumber string, newAccount *Account) (*Account, error) {
	tx := configs.DB.Begin()

	if err := tx.Model(&Account{}).Where("account_number = ?", accNumber).Update("nominal", newAccount.Nominal).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return newAccount, nil
}
