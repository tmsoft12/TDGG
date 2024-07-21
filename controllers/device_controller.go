package controllers

import (
	"context"
	"tm/database"
	"tm/models"

	"github.com/gofiber/fiber/v2"
)

func CreateDoorLockData(c *fiber.Ctx) error {
	data := new(models.DoorLockData)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	sql := `
		INSERT INTO door_lock_data (
			device_id, latitude, longitude, battery_level, signal_status, is_locked
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
	`
	_, err := database.DBpool.Exec(context.Background(), sql,
		data.DeviceID,
		data.Latitude,
		data.Longitude,
		data.BatteryLevel,
		data.SignalStatus,
		data.IsLocked,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(data)
}

// Read
func GetDoorLockData(c *fiber.Ctx) error {
	deviceID := c.Params("device_id")

	var data models.DoorLockData
	sql := `
		SELECT device_id, latitude, longitude, battery_level, signal_status, is_locked
		FROM door_lock_data
		WHERE device_id = $1
	`
	row := database.DBpool.QueryRow(context.Background(), sql, deviceID)
	err := row.Scan(
		&data.DeviceID,
		&data.Latitude,
		&data.Longitude,
		&data.BatteryLevel,
		&data.SignalStatus,
		&data.IsLocked,
	)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}

	return c.JSON(data)
}

// Update
func UpdateDoorLockData(c *fiber.Ctx) error {
	data := new(models.DoorLockData)
	if err := c.BodyParser(data); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	sql := `
		UPDATE door_lock_data
		SET latitude = $1, longitude = $2, battery_level = $3, signal_status = $4, is_locked = $5
		WHERE device_id = $6
	`
	_, err := database.DBpool.Exec(context.Background(), sql,
		data.Latitude,
		data.Longitude,
		data.BatteryLevel,
		data.SignalStatus,
		data.IsLocked,
		data.DeviceID,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(data)
}

// Delete
func DeleteDoorLockData(c *fiber.Ctx) error {
	deviceID := c.Params("device_id")

	sql := `
		DELETE FROM door_lock_data
		WHERE device_id = $1
	`
	_, err := database.DBpool.Exec(context.Background(), sql, deviceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
