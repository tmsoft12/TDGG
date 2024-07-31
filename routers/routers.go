package routes

import (
	"tm/controllers"
	"tm/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutes(app *fiber.App) {
	// CORS middleware
	app.Use(middlewares.CORSMiddleware())

	// Swagger UI
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Auth
	app.Post("/api/login", controllers.Login)

	// User routes
	userGroup := app.Group("/api", middlewares.OnlyUser)

	// Device routes
	userGroup.Get("/device/all_device", controllers.GetAllDevices)
	userGroup.Get("/device/last_locations", controllers.GetAllDevicesLastLocation)
	userGroup.Post("/device/locations", controllers.AddDeviceLocation)
	userGroup.Get("/device/location_list/:id", controllers.GetDeviceLocations)

	// Driver routes
	userGroup.Get("/driver/all_driver", controllers.GetAllDrivers)
	userGroup.Get("/driver/get_driver/:id", controllers.GetDriverById)
	userGroup.Post("/driver/create_driver", controllers.CreateDriver)
	userGroup.Delete("/driver/delete_driver/:id", controllers.DeleteDriver)
	userGroup.Put("/driver/update_driver/:id", controllers.UpdateDriver)

	// Home page route
	userGroup.Get("/main", controllers.Home_page)

	// Admin routes
	adminGroup := app.Group("/api/admin", middlewares.OnlyAdmin)
	adminGroup.Get("/allusers", controllers.GetAllUser)
	adminGroup.Post("/createuser", controllers.CreateUser)
	adminGroup.Get("/getuser/:id", controllers.GetUserById)
	adminGroup.Put("/update/:id", controllers.UpdateUser)
	adminGroup.Delete("/delete/:id", controllers.DeleteUser)

	// WebSocket route
	app.Get("/socket", websocket.New(controllers.HandleConnection))
}
