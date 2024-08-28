package routes

import (
	"account_services/src/controllers"
	"account_services/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Server is running",
		})
	})
	app.Get("/users", controllers.FindAllUsers)
	app.Post("/login", controllers.Login)


	app.Get("/saldo/:no_rekening", controllers.GetSaldo)
	app.Post("/daftar", controllers.Register)
	app.Post("/tabung", controllers.DepositBalance)
	app.Post("/tarik", middlewares.JwtMiddleware(), controllers.WithdrawalBalance)
}
