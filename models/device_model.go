package models

type DoorLockData struct {
	DeviceID     string   `json:"device_id"`
	Latitude     *float64 `json:"latitude,omitempty"`
	Longitude    *float64 `json:"longitude,omitempty"`
	BatteryLevel int      `json:"battery_level"`
	SignalStatus string   `json:"signal_status"`
	IsLocked     bool     `json:"is_locked"`
}
