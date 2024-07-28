package controllers

import (
	"context"
	"log"
	"time"
	"tm/database"
	"tm/models"

	"github.com/gofiber/fiber/v2"
)

// @Summary Get all devices with their last known location
// @Description Get all devices with their last known location
// @Tags Devices
// @Produce json
// @Success 200 {array} models.DeviceSchema
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/device/last_locations [get]
func GetAllDevicesLastLocation(c *fiber.Ctx) error {

	rows, err := database.DBpool.Query(context.Background(), "SELECT device_id, battery_level, signal_status, is_locked,status FROM devices")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving devices"})
	}
	defer rows.Close()

	var devices []models.DeviceSchema
	for rows.Next() {
		var device models.DeviceSchema
		err := rows.Scan(&device.DeviceId, &device.BatteryLevel, &device.SignalStatus, &device.IsLocked, &device.Status)
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

// @Summary Add device location
// @Description Add a new location entry for a device
// @Tags Devices
// @Accept json
// @Produce json
// @Param location body models.DeviceLocationRequest true "Device location"
// @Success 201 {object} map[string]interface{} "Location added successfully"
// @Failure 400 {object} map[string]interface{} "Bad request"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/device/locations [post]
func AddDeviceLocation(c *fiber.Ctx) error {
	// Request gövdesinden JSON verisini alın
	var requestData models.DeviceLocationRequest

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

// @Summary Get device locations
// @Description Get all locations for a specific device
// @Tags Devices
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {array} models.DeviceLocation
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/device/location_list/{id} [get]
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

// @Summary Get all devices
// @Description Get all devices
// @Tags Devices
// @Produce json
// @Success 200 {array} models.DeviceAll
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /api/device/all_device [get]
func GetAllDevices(c *fiber.Ctx) error {

	rows, err := database.DBpool.Query(context.Background(), "SELECT device_id, battery_level, signal_status, is_locked,status FROM devices")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving devices"})
	}
	defer rows.Close()

	var devices []models.DeviceAll
	for rows.Next() {
		var device models.DeviceAll
		err := rows.Scan(&device.DeviceId, &device.BatteryLevel, &device.SignalStatus, &device.IsLocked, &device.Status)
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

// Home_page retrieves the count of devices by their status and the latest location of each device
// @Summary Get device status counts and latest locations
// @Description Retrieves the count of devices grouped by their status and the latest location of each device
// @Tags devices
// @Accept json
// @Produce json
// @Success 200 {object} HomePageResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/main [get]
func Home_page(c *fiber.Ctx) error {
	// Query to get the count of devices grouped by status
	query := `SELECT status, COUNT(*) as count FROM devices GROUP BY status;`
	rows, err := database.DBpool.Query(context.Background(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get status counts",
		})
	}
	defer rows.Close()

	var statusCounts []models.StatusCount
	for rows.Next() {
		var sc models.StatusCount
		err := rows.Scan(&sc.Status, &sc.Count)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse status count",
			})
		}
		statusCounts = append(statusCounts, sc)
	}

	if rows.Err() != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error occurred during rows iteration",
		})
	}

	// Query to get the latest location of each device
	query = "SELECT device_id, battery_level, signal_status, is_locked, status FROM devices"
	rows, err = database.DBpool.Query(context.Background(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error retrieving devices"})
	}
	defer rows.Close()

	var devices []models.DeviceSchema
	for rows.Next() {
		var device models.DeviceSchema
		err := rows.Scan(&device.DeviceId, &device.BatteryLevel, &device.SignalStatus, &device.IsLocked, &device.Status)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error scanning device"})
		}

		// Query to get the latest location of the device
		locationQuery := "SELECT timestamp, latitude, longitude FROM device_locations WHERE device_id=$1 ORDER BY timestamp DESC LIMIT 1"
		locationRow := database.DBpool.QueryRow(context.Background(), locationQuery, device.DeviceId)
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"Status":    statusCounts,
		"Locations": devices,
	})
}

// HomePageResponse represents the response structure for the Home_page endpoint
type HomePageResponse struct {
	Status        []models.StatusCount  `json:"Status"`
	LastLocations []models.DeviceSchema `json:"Last Locations"`
}

// ErrorResponse represents the error response structure
type ErrorResponse struct {
	Error string `json:"error"`
}
