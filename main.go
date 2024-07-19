package main

import (
	"tm/database"
	routes "tm/routers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Veritabanı bağlantısını başlat
	database.InitDB()

	routes.SetupRoutes(app)

	app.Listen(":8000")
}
