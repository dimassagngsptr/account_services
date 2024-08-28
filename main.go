package main

import (
	"account_services/src/configs"
	"account_services/src/helpers"
	"account_services/src/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	HOST := os.Getenv("API_HOST")
	PORT := os.Getenv("API_PORT")
	helpers.InitLogger()

	app := fiber.New()
	configs.InitDB()
	helpers.Migration()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PATCH, DELETE",
	}))

	routes.Routes(app)
	helpers.LogWithFields(logrus.InfoLevel, "Application", "INIT", nil, "Application started")

	log.Fatal(app.Listen(HOST + ":" + PORT))
}
