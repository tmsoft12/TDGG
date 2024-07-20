package routes

import (
	"tm/controllers"
	"tm/middlewares"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,OPTIONS",
	}))

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
