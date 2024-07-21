package main

import (
	"tm/database"
	routes "tm/routers"

	_ "tm/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// @title Fiber Swagger Example API
// @version 1.0
// @description This is a sample server for a Fiber API.
// @host localhost:8000
// @BasePath /

func main() {
	database.InitDB()
	app := fiber.New()
	app.Use(logger.New())

	routes.SetupRoutes(app)

	// default

	app.Listen("0.0.0.0:8000")

}
