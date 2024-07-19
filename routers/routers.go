package routes

import (
	"tm/controllers"
	"tm/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/login", controllers.Login)

	userGroup := app.Group("/home", middlewares.OnlyUser)
	userGroup.Get("/", controllers.UserTest)

	adminGroup := app.Group("/admin", middlewares.OnlyAdmin)
	adminGroup.Get("/allusers", controllers.GetAllUser)
	adminGroup.Post("/createuser", controllers.CreateUser)
	adminGroup.Get("/getuser/:id", controllers.GetUserById)
	adminGroup.Put("/update/:id", controllers.UpdateUser)
	adminGroup.Delete("/delete/:id", controllers.DeleteReport)

}
