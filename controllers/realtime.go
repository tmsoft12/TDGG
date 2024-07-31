package controllers

import (
	"context"
	"encoding/json"
	"log"
	"tm/database"
	"tm/models"

	"github.com/gofiber/websocket/v2"
)

// WebSocket upgrader ve client map
var clients = make(map[*websocket.Conn]bool)

// WebSocket bağlantılarını yönetir
func HandleConnection(c *websocket.Conn) {
	defer c.Close()

	clients[c] = true
	log.Println("Client connected")

	// Eski verileri gönder
	if err := sendInitialData(c); err != nil {
		log.Println("Error sending initial data:", err)
		c.Close()
		delete(clients, c)
		return
	}

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			log.Println("Error while reading message:", err)
			delete(clients, c)
			break
		}
	}
}

// Eski verileri alır ve WebSocket istemcisine gönderir
func sendInitialData(c *websocket.Conn) error {
	data, err := fetchAllData()
	if err != nil {
		return err
	}

	message, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return c.WriteMessage(websocket.TextMessage, message)
}

// Veritabanından eski verileri alır
func fetchAllData() ([]models.DeviceAll, error) {
	rows, err := database.DBpool.Query(context.Background(), "SELECT device_id, battery_level, signal_status, is_locked, status FROM devices")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []models.DeviceAll
	for rows.Next() {
		var device models.DeviceAll
		if err := rows.Scan(&device.DeviceId, &device.BatteryLevel, &device.SignalStatus, &device.IsLocked, &device.Status); err != nil {
			return nil, err
		}

		// Son konum bilgisini al
		locationRow := database.DBpool.QueryRow(context.Background(), "SELECT timestamp, latitude, longitude FROM device_locations WHERE device_id=$1 ORDER BY timestamp DESC LIMIT 1", device.DeviceId)
		var location models.DeviceLocation
		err = locationRow.Scan(&location.Timestamp, &location.Latitude, &location.Longitude)
		if err != nil {
			device.Location = nil // Koordinat bilgisi yok
		} else {
			device.Location = &location
		}

		devices = append(devices, device)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

// Tüm bağlı istemcilere güncellemeleri yayar
func broadcastUpdate(data []models.DeviceAll) {
	message, err := json.Marshal(data)
	if err != nil {
		log.Println("Error while marshaling message:", err)
		return
	}

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error while writing message:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// Veritabanından gelen güncellemeleri dinler
func ListenForUpdates() {
	conn, err := database.DBpool.Acquire(context.Background())
	if err != nil {
		log.Fatal("Unable to acquire database connection:", err)
	}
	defer conn.Release()

	_, err = conn.Conn().Exec(context.Background(), "LISTEN data_update")
	if err != nil {
		log.Fatal("Unable to listen for notifications:", err)
	}

	log.Println("Listening for updates")

	for {
		notification, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			log.Println("Error while waiting for notification:", err)
			continue
		}

		log.Printf("Received notification: %s", notification.Payload)

		// Güncellenmiş verileri al ve yayınla
		data, err := fetchAllData()
		if err != nil {
			log.Println("Error fetching updated data:", err)
			continue
		}
		broadcastUpdate(data)
	}
}
