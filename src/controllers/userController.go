package controllers

import (
	"account_services/src/configs"
	"account_services/src/helpers"
	"account_services/src/middlewares"
	"account_services/src/models"
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func FindAllUsers(c *fiber.Ctx) error {
	users := models.FindAllUsers()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": users,
	})
}

func Register(c *fiber.Ctx) error {
	var payload models.User

	if err := c.BodyParser(&payload); err != nil {
		helpers.LogWithFields(logrus.ErrorLevel, "Register", c.Method(), payload, fmt.Sprintf("Failed to parse request body %v", err.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"remark": "Invalid request body",
		})
	}

	helpers.LogWithFields(logrus.InfoLevel, "Register", c.Method(), payload, "Received request to create user")

	errors := helpers.Validate(payload)
	if len(errors) > 0 {
		var errorMessages []string
		for _, err := range errors {
			errorMessages = append(errorMessages, fmt.Sprintf("%v", err))
		}
		helpers.LogWithFields(logrus.WarnLevel, "Register", c.Method(), payload, fmt.Sprintf("Failed to parse request body: %v", strings.Join(errorMessages, ", ")))
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"remark": "Failed to register request",
			"error":  errorMessages,
		})
	}

	existingUser := models.FindByNikorPhone(payload.NIK, payload.Phone)
	if existingUser != nil {
		helpers.LogWithFields(logrus.WarnLevel, "Register", "GET", payload, "Account already exists")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"remark": "Account already exists",
		})
	}
	newUser, err := models.CreateUser(&payload)
	if err != nil {
		helpers.LogWithFields(logrus.ErrorLevel, "Register", c.Method(), payload, fmt.Sprintf("Failed to create new user %v", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"remark": "Internal Server Error",
		})
	}
	newAccount := models.Account{
		AccountNumber: fmt.Sprintf("70001%s", payload.Phone),
		UserID:        newUser.ID.String(),
		Nominal:       0,
	}
	account, err := models.CreateAccount(&newAccount)
	if err != nil {
		helpers.LogWithFields(logrus.ErrorLevel, "Register", c.Method(), newAccount, fmt.Sprintf("Failed to create account user %v", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"remark": "Internal server error",
		})
	}

	logger := map[string]interface{}{
		"user":    *newUser,
		"account": *account,
	}

	helpers.LogWithFields(logrus.InfoLevel, "Register", c.Method(), logger, "Success to create user and account")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"no_rekening": account.AccountNumber,
	})

}

func Login(c *fiber.Ctx) error {
	var input models.LoginRequest
	if err := c.BodyParser(&input); err != nil {
		helpers.LogWithFields(logrus.ErrorLevel, "Login", c.Method(), input, fmt.Sprintf("Failed to parse request body %v", err.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"remark": "Invalid request body",
		})
	}
	var user models.User
	result := configs.DB.Preload("Account", func(db *gorm.DB) *gorm.DB {
		var account []*models.APIAccount
		return db.Model(&models.Account{}).Find(&account)
	}).First(&user, "phone = ?", input.Phone)
	if result.RowsAffected == 0 {
		helpers.LogWithFields(logrus.WarnLevel, "Login", "GET", input.Phone, "User not found")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"remark": "User not found",
		})
	}

	helpers.LogWithFields(logrus.InfoLevel, "Login", "GET", user, "User found")

	var accountNumber string
	if len(user.Account) > 0 {
		for _, account := range user.Account {
			accountNumber = account.AccountNumber
		}
	}

	jwtKey := os.Getenv("SECRETKEY")
	payload := map[string]interface{}{
		"ID":             user.ID,
		"account_number": accountNumber,
	}
	token, err := middlewares.GenerateToken(jwtKey, payload)
	if err != nil {
		helpers.LogWithFields(logrus.InfoLevel, "Login", "GET", token, fmt.Sprintf("Failed to generate token: %v", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"remark": "Failed to generate Token",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token":   token,
		"data":    user,
	})
}
