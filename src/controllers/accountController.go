package controllers

import (
	"account_services/src/helpers"
	"account_services/src/models"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func DepositBalance(c *fiber.Ctx) error {
	var payload models.AccountRequest

	if err := c.BodyParser(&payload); err != nil {
		helpers.LogWithFields(logrus.ErrorLevel, "DepositAccount", c.Method(), payload, fmt.Sprintf("Failed to parse request body %v", err.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"remark": "Invalid request body",
		})
	}

	helpers.LogWithFields(logrus.InfoLevel, "DepositAccount", c.Method(), payload, "Received request to deposit account")

	if payload.Nominal < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"remark": "Nominal must be greater than 0",
		})
	}
	existingAccount := models.FindByAccNumber(payload.AccountNumber)
	if existingAccount == nil {
		helpers.LogWithFields(logrus.WarnLevel, "DepositAccount", "GET", payload.AccountNumber, "Account not found during deposit")
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"remark": "Account not found",
		})
	}
	newSaldo := existingAccount.Nominal + payload.Nominal
	existingAccount.Nominal = newSaldo

	result, err := models.UpdateBalance(existingAccount.AccountNumber, existingAccount)
	if err != nil {
		helpers.LogWithFields(logrus.ErrorLevel, "DepositAccount", "UPDATE", payload, fmt.Sprintf("Failed to update account balance %v", err.Error()))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"remark": "Internal server error",
		})
	}

	helpers.LogWithFields(logrus.InfoLevel, "DepositAccount", "UPDATE", existingAccount.Nominal, "Account balance updated successfully")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"saldo": result.Nominal,
	})
}

func WithdrawalBalance(c *fiber.Ctx) error {
	var payload models.AccountRequest

	if err := c.BodyParser(&payload); err != nil {
		helpers.LogWithFields(logrus.ErrorLevel, "WithdrawalAccount", c.Method(), payload, fmt.Sprintf("Failed to parse request body %v", err.Error()))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"remark": "Invalid request body",
		})
	}

	user, ok := c.Locals("user").(jwt.MapClaims)
	accNumber := user["account_number"]
	if !ok || accNumber != payload.AccountNumber {
		claims := map[string]interface{}{
			"account_number": accNumber,
			"payload":        payload.AccountNumber,
		}
		helpers.LogWithFields(logrus.WarnLevel, "WithdrawalAccount", c.Method(), claims, "Account number not same")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"remark": "Unauthorized",
		})
	}
	helpers.LogWithFields(logrus.InfoLevel, "WithdrawalAccount", c.Method(), user, "Received token claims")

	helpers.LogWithFields(logrus.InfoLevel, "WithdrawalAccount", c.Method(), payload, "Received request to withdrawal account")

	existingAccount := models.FindByAccNumber(payload.AccountNumber)
	if existingAccount == nil {

		helpers.LogWithFields(logrus.WarnLevel, "WithdrawalAccount", "GET", payload.AccountNumber, "Account not found withdrawal deposit")

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"remark": "Account not found",
		})
	}
	if existingAccount.Nominal < payload.Nominal || existingAccount.Nominal == 0 {
		helpers.LogWithFields(logrus.WarnLevel, "WithdrawalAccount", "GET", existingAccount.Nominal, "Your balance is insufficient")

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"remark": "Your balance is insufficient",
		})
	}

	newSaldo := existingAccount.Nominal - payload.Nominal
	existingAccount.Nominal = newSaldo

	result, err := models.UpdateBalance(existingAccount.AccountNumber, existingAccount)
	if err != nil {
		helpers.LogWithFields(logrus.ErrorLevel, "WithdrawalAccount", "UPDATE", payload, fmt.Sprintf("Internal server error failed updating balance %v", err.Error()))

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"remark": "Internal server error",
		})
	}

	helpers.LogWithFields(logrus.InfoLevel, "WithdrawalAccount", c.Method(), result, "Account balance updated successfully")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"saldo": result.Nominal,
	})
}

func GetSaldo(c *fiber.Ctx) error {
	accountNumber := c.Params("no_rekening")
	foundAccount := models.FindByAccNumber(accountNumber)
	if foundAccount == nil {

		helpers.LogWithFields(logrus.WarnLevel, "GetAccount", "GET", accountNumber, fmt.Sprintf("Account not found with account number %v", accountNumber))

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"remark": "Account not found",
		})
	}
	helpers.LogWithFields(logrus.InfoLevel, "GetAccount", c.Method(), foundAccount, fmt.Sprintf("Success get account with acc number %v", accountNumber))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"saldo": foundAccount.Nominal,
	})
}
