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
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// @Summary Login
// @Description Login
// @Tags User
// @Accept json
// @Produce json
// @Param loginInput body LoginInput true "Login Input"
// @Success 200 {object} map[string]interface{} "token"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 401 {object} map[string]interface{} "Invalid username or password"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/login [post]
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

	return c.JSON(fiber.Map{"token": tokenString, "Rule": user.Role})
}

// @Summary Get All Users
// @Description Retrieve all users
// @Tags Admin
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/allusers [get]
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

// @Summary Get User By ID
// @Description Retrieve a user by ID
// @Tags Admin
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/getuser/{id} [get]
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

// @Summary Create User
// @Description Create a new user
// @Tags Admin
// @Accept json
// @Produce json
// @Param user body models.User true "User"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/createuser [post]
func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	if user.Username == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username and password cannot be empty"})
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

// @Summary Update User
// @Description Update an existing user
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.User true "User"
// @Success 200 {object} map[string]interface{} "User updated"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/update/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body: " + err.Error())
	}

	if user.Username == "" || user.Password == "" || user.Role == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Username, password, and role must not be empty")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error hashing password")
	}

	commandTag, err := database.DBpool.Exec(context.Background(),
		"UPDATE users SET username=$1, password=$2, role=$3 WHERE id=$4",
		user.Username, hashedPassword, user.Role, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Database error: " + err.Error())
	}

	if commandTag.RowsAffected() == 0 {
		return c.Status(fiber.StatusNotFound).SendString("User not found or no changes were made")
	}

	return c.SendStatus(fiber.StatusOK)
}

// @Summary Delete User
// @Description Delete a user by ID
// @Tags Admin
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User deleted"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/admin/delete/{id} [delete]
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
