package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v4"
)

var JwtSecret = []byte("your_secret_key")

func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	})
}

func parseToken(c *fiber.Ctx) (*jwt.Token, error) {
	// Check for token in cookies
	cookie := c.Cookies("jwt")
	if cookie != "" {
		return jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}
			return JwtSecret, nil
		})
	}

	// Check for token in Authorization header
	authHeader := c.Get("Authorization")
	if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenStr := authHeader[7:]
		return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}
			return JwtSecret, nil
		})
	}

	return nil, fiber.NewError(fiber.StatusUnauthorized, "unauthenticated")
}

func OnlyUser(c *fiber.Ctx) error {
	token, err := parseToken(c)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "You are not login"})
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["role"] != "user" && claims["role"] != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "You are not login"})
	}

	return c.Next()
}

func OnlyAdmin(c *fiber.Ctx) error {
	token, err := parseToken(c)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "unauthenticated"})
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["role"] != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "You are not admin"})
	}

	return c.Next()
}
