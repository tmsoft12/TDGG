package main

import (
	"tm/database"
	routes "tm/routers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	database.InitDB()
	app := fiber.New()
	app.Use(logger.New())

	routes.SetupRoutes(app)

	// default

	app.Listen("0.0.0.0:8000")

}
