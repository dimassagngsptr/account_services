package models

import (
	"account_services/src/configs"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID      uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	Name    string       `json:"name" validate:"required"`
	NIK     string       `json:"nik" validate:"required"`
	Phone   string       `json:"phone" validate:"required,min=10,max=13"`
	Account []APIAccount `json:"account"`
}

type APIAccount struct {
	AccountNumber string  `json:"account_number"`
	Nominal       float64 `json:"nominal"`
	UserID        string  `json:"user_id"`
}

type LoginRequest struct {
	Phone string `json:"phone"`
}

func FindAllUsers() []*User {
	var users []*User
	configs.DB.Preload("Account", func(db *gorm.DB) *gorm.DB {
		var account []*APIAccount
		return db.Model(&Account{}).Find(&account)
	}).Find(&users)
	return users
}

func CreateUser(newUser *User) (*User, error) {
	if err := configs.DB.Create(&newUser).Error; err != nil {
		return nil, err
	}
	return newUser, nil
}

func FindByNikorPhone(nik, phone string) *User {
	var user User
	result := configs.DB.Where("nik = ? OR phone = ?", nik, phone).First(&user)
	if result.RowsAffected == 0 {
		return nil
	}
	return &user
}
