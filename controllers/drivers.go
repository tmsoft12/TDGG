package controllers

import (
	"context"
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
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
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
