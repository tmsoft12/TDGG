package routes

import (
	"tm/controllers"
	"tm/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(middlewares.CORSMiddleware())

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Post("api/login", controllers.Login)

	//USER////////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////
	userGroup := app.Group("/api", middlewares.OnlyUser)
	//DEVICE//////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////
	userGroup.Get("/device/all_device", controllers.GetAllDevices)
	userGroup.Get("/device/last_locations", controllers.GetAllDevicesLastLocation)
	userGroup.Post("/device/locations", controllers.AddDeviceLocation)
	userGroup.Get("/device/location_list/:id", controllers.GetDeviceLocations)

	//////////////////////////////////////////////////////////////////////////
	//ADMIN///////////////////////////////////////////////////////////////////
	//////////////////////////////////////////////////////////////////////////
	adminGroup := app.Group("api/admin", middlewares.OnlyAdmin)
	adminGroup.Get("/allusers", controllers.GetAllUser)
	adminGroup.Post("/createuser", controllers.CreateUser)
	adminGroup.Get("/getuser/:id", controllers.GetUserById)
	adminGroup.Put("/update/:id", controllers.UpdateUser)
	adminGroup.Delete("/delete/:id", controllers.DeleteUser)

	//////////////////////////////////////////////////////////////////////////

}
