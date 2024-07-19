package controllers

import (
	"context"
	"time"
	"tm/database"
	"tm/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("your_secret_key")

type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Login godoc
// @Summary User login
// @Description Login by username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param loginInput body LoginInput true "Login input"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /admin/login [post]
func Login(c *fiber.Ctx) error {
	input := new(LoginInput)
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var user models.User
	err := database.DBpool.QueryRow(context.Background(), "SELECT id, username, password, role FROM users WHERE username=$1", input.Username).Scan(&user.Id, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	if !checkPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
	})

	return c.JSON(fiber.Map{"token": tokenString})
}

// GetAllUser godoc
// @Summary Get all users
// @Description Get details of all users
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {string} string "Internal Server Error"
// @Router /admin/users [get]
func GetAllUser(c *fiber.Ctx) error {
	rows, err := database.DBpool.Query(context.Background(), "SELECT id, username, password, role FROM users")
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Role)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(fiber.Map{"users": users})
}

// GetUserById godoc
// @Summary Get user by ID
// @Description Get details of a user by ID
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /admin/users/{id} [get]
func GetUserById(c *fiber.Ctx) error {
	id := c.Params("id")

	row := database.DBpool.QueryRow(context.Background(), "SELECT id, username, password, role FROM users WHERE id = $1", id)

	var user models.User
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString("User not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(fiber.Map{"user": user})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the given details
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User input"
// @Success 201 {object} models.User
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /admin/users [post]
func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(500).SendString("Error hashing password")
	}

	err = database.DBpool.QueryRow(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", user.Username, hashedPassword).Scan(&user.Id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(201).JSON(user)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update details of an existing user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true "User input"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /admin/users/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	requestBody := new(struct {
		User struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Password string `json:"password"`
			Role     string `json:"role"`
		} `json:"user"`
	})

	if err := c.BodyParser(requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body: " + err.Error())
	}

	user := requestBody.User

	if user.Username == "" || user.Password == "" || user.Role == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Username, password, and role must not be empty")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
	}

	commandTag, err := database.DBpool.Exec(context.Background(),
		"UPDATE users SET username=$1, password=$2, role=$3 WHERE id=$4",
		user.Username, hashedPassword, user.Role, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Database error: " + err.Error())
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).SendString("User not found or no changes were made")
	}

	return c.SendStatus(fiber.StatusOK)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Param id path string true "User ID"
// @Success 200 {string} string "OK"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /admin/users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	commandTag, err := database.DBpool.Exec(context.Background(), "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	if commandTag.RowsAffected() == 0 {
		return c.Status(404).SendString("User not found")
	}

	return c.SendStatus(200)
}
