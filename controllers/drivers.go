package controllers

import (
	"context"
	"strconv"
	"time"
	"tm/database"
	"tm/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
)

// @Summary Get all drivers
// @Description Get all devices
// @Tags Drivers
// @Produce json
// @Success 200 {array} models.Driver
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/driver/all_driver [get]
func GetAllDrivers(c *fiber.Ctx) error {
	rows, err := database.DBpool.Query(context.Background(), "SELECT * FROM driver")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}
	defer rows.Close()
	var drivers []models.Driver
	for rows.Next() {
		var driver models.Driver
		err := rows.Scan(&driver.ID, &driver.CreateTime, &driver.Name, &driver.Phone, &driver.CarNumber, &driver.CarModel, &driver.Weight, &driver.Country)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
		}
		drivers = append(drivers, driver)
	}

	return c.Status(fiber.StatusOK).JSON(drivers)
}

// @Summary Get Driver By ID
// @Description Retrieve a driver by ID
// @Tags Drivers
// @Produce json
// @Param id path int true "Driver ID"
// @Success 200 {object} models.Driver
// @Failure 404 {object} map[string]interface{} "Driver not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/driver/get_driver/{id} [get]
func GetDriverById(c *fiber.Ctx) error {
	id := c.Params("id")
	row := database.DBpool.QueryRow(context.Background(), "SELECT * FROM driver WHERE id = $1", id)
	var driver models.Driver
	err := row.Scan(&driver.ID, &driver.CreateTime, &driver.Name, &driver.Phone, &driver.CarNumber, &driver.CarModel, &driver.Weight, &driver.Country)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
	}
	return c.JSON(fiber.Map{
		"driver": driver,
	})
}

// @Summary Create Driver
// @Description Create a new driver
// @Tags Drivers
// @Accept json
// @Produce json
// @Param driver body models.Driver true "Driver"
// @Success 201 {object} models.Driver
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/driver/create_driver [post]
func CreateDriver(c *fiber.Ctx) error {
	// Parse the request body into the driver model
	driver := new(models.Driver)
	if err := c.BodyParser(driver); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "cannot parse JSON",
			"message": err.Error(),
		})
	}

	// Set the creation time
	driver.CreateTime = time.Now()

	// Insert the new driver into the database
	query := `
        INSERT INTO driver (create_time, name, phone, car_number, car_model, weight, country) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) 
        RETURNING id`
	var driverID int
	err := database.DBpool.QueryRow(
		context.Background(),
		query,
		driver.CreateTime,
		driver.Name,
		driver.Phone,
		driver.CarNumber,
		driver.CarModel,
		driver.Weight,
		driver.Country).Scan(&driverID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to insert driver into database",
			"message": err.Error(),
		})
	}

	// Set the new driver's ID
	driver.ID = driverID

	// Return the newly created driver
	return c.Status(fiber.StatusCreated).JSON(driver)
}

// @Summary Delete Driver
// @Description Delete a driver by ID
// @Tags Drivers
// @Accept json
// @Produce json
// @Param id path int true "Driver ID"
// @Success 200 {object} map[string]interface{} "Deleted successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/driver/delete_driver/{id} [delete]
func DeleteDriver(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID parameter is required",
		})
	}

	query := `DELETE FROM driver WHERE id = $1`
	result, err := database.DBpool.Exec(context.Background(), query, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to delete driver from database",
			"message": err.Error(),
		})
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "driver not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Deleted successfully",
	})
}

// @Summary Update Driver
// @Description Update a driver's details by ID
// @Tags Drivers
// @Accept json
// @Produce json
// @Param id path int true "Driver ID"
// @Param driver body models.Driver true "Driver"
// @Success 200 {object} models.Driver
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 404 {object} map[string]interface{} "Not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/driver/update_driver/{id} [put]
func UpdateDriver(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID parameter is required",
		})
	}

	// Parse the request body into the driver model
	driver := new(models.Driver)
	if err := c.BodyParser(driver); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "cannot parse JSON",
			"message": err.Error(),
		})
	}

	// Set the ID from the URL parameter
	driver.ID, _ = strconv.Atoi(id)

	// Update the driver's details in the database
	query := `
        UPDATE driver 
        SET name = $1, phone = $2, car_number = $3, car_model = $4, weight = $5, country = $6
        WHERE id = $7`
	result, err := database.DBpool.Exec(
		context.Background(),
		query,
		driver.Name,
		driver.Phone,
		driver.CarNumber,
		driver.CarModel,
		driver.Weight,
		driver.Country,
		driver.ID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "failed to update driver in database",
			"message": err.Error(),
		})
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "driver not found",
		})
	}

	// Return the updated driver
	return c.Status(fiber.StatusOK).JSON(driver)
}
