package controllers

import (
	"context"
	"log"
	"time"
	"tm/database"
	"tm/models"

	"github.com/gofiber/fiber/v2"
)

func GetAllDevicesLastLocation(c *fiber.Ctx) error {

	rows, err := database.DBpool.Query(context.Background(), "SELECT device_id, battery_level, signal_status, is_locked FROM devices")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving devices"})
	}
	defer rows.Close()

	var devices []models.DeviceSchema
	for rows.Next() {
		var device models.DeviceSchema
		err := rows.Scan(&device.DeviceId, &device.BatteryLevel, &device.SignalStatus, &device.IsLocked)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning device"})
		}

		locationRow := database.DBpool.QueryRow(context.Background(), "SELECT timestamp, latitude, longitude FROM device_locations WHERE device_id=$1 ORDER BY timestamp DESC LIMIT 1", device.DeviceId)
		var location models.DeviceLocation
		err = locationRow.Scan(&location.Timestamp, &location.Latitude, &location.Longitude)
		if err != nil {

			device.Location = models.DeviceLocation{}
		} else {
			device.Location = location
		}

		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing devices"})
	}

	return c.Status(fiber.StatusOK).JSON(devices)
}

func AddDeviceLocation(c *fiber.Ctx) error {
	// Request gövdesinden JSON verisini alın
	var requestData struct {
		DeviceId  string  `json:"device_id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse request", "details": err.Error()})
	}

	// deviceId'yi JSON verisinden alın
	deviceId := requestData.DeviceId

	// deviceId'yi ekrana yazdır
	log.Printf("Received deviceId: %s\n", deviceId)

	// Güncel Unix zaman damgasını al
	timestamp := time.Now().Unix()

	// Veritabanına yeni location bilgisini ekleyin
	_, err := database.DBpool.Exec(
		context.Background(),
		"INSERT INTO device_locations (device_id, timestamp, latitude, longitude) VALUES ($1, $2, $3, $4)",
		deviceId, timestamp, requestData.Latitude, requestData.Longitude,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add location", "details": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Location added successfully"})
}

func GetDeviceLocations(c *fiber.Ctx) error {
	deviceId := c.Params("id")

	rows, err := database.DBpool.Query(context.Background(), "SELECT timestamp, latitude, longitude FROM device_locations WHERE device_id=$1 ORDER BY timestamp DESC", deviceId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving locations"})
	}
	defer rows.Close()

	var locations []models.DeviceLocation
	for rows.Next() {
		var location models.DeviceLocation
		err := rows.Scan(&location.Timestamp, &location.Latitude, &location.Longitude)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning location"})
		}
		locations = append(locations, location)
	}

	if err = rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing locations"})
	}

	return c.Status(fiber.StatusOK).JSON(locations)
}

func GetAllDevices(c *fiber.Ctx) error {

	rows, err := database.DBpool.Query(context.Background(), "SELECT device_id, battery_level, signal_status, is_locked FROM devices")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving devices"})
	}
	defer rows.Close()

	var devices []models.DeviceAll
	for rows.Next() {
		var device models.DeviceAll
		err := rows.Scan(&device.DeviceId, &device.BatteryLevel, &device.SignalStatus, &device.IsLocked)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning device"})
		}

		devices = append(devices, device)
	}

	if err = rows.Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error processing devices"})
	}

	return c.Status(fiber.StatusOK).JSON(devices)
}
