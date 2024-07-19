package controllers

import "github.com/gofiber/fiber/v2"

func UserTest(c *fiber.Ctx) error {
	return c.JSON("Home page hos geldin")
}
