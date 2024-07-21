package routes

import (
	"tm/controllers"
	"tm/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func SetupRoutes(app *fiber.App) {
	app.Use(middlewares.CORSMiddleware())
	app.Post("/api/doorlock", controllers.CreateDoorLockData)
	app.Get("/api/doorlock/:device_id", controllers.GetDoorLockData)
	app.Put("/api/doorlock/:device_id", controllers.UpdateDoorLockData)
	app.Delete("/api/doorlock/:device_id", controllers.DeleteDoorLockData)

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Post("api/login", controllers.Login)
	//USER
	userGroup := app.Group("/api", middlewares.OnlyUser)
	userGroup.Get("/home", controllers.UserTest)

	//ADMIN
	adminGroup := app.Group("api/admin", middlewares.OnlyAdmin)
	adminGroup.Get("/allusers", controllers.GetAllUser)
	adminGroup.Post("/createuser", controllers.CreateUser)
	adminGroup.Get("/getuser/:id", controllers.GetUserById)
	adminGroup.Put("/update/:id", controllers.UpdateUser)
	adminGroup.Delete("/delete/:id", controllers.DeleteUser)

}
