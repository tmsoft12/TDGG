package models

type DeviceLocation struct {
	Timestamp int64   `json:"timestamp"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type SingleDeviceSchema struct {
	DeviceId     string           `json:"deviceId"`
	Location     []DeviceLocation `json:"location"`
	BatteryLevel int              `json:"batteryLevel"`
	SignalStatus string           `json:"signalStatus"`
	IsLocked     bool             `json:"isLocked"`
}

type DeviceSchema struct {
	DeviceId     string         `json:"deviceId"`
	Location     DeviceLocation `json:"location"`
	BatteryLevel int            `json:"batteryLevel"`
	SignalStatus string         `json:"signalStatus"`
	IsLocked     bool           `json:"isLocked"`
	Status       string         `json:"status"`
}

// DeviceAll yapısı

type DeviceLocationRequest struct {
	DeviceId  string  `json:"device_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type DeviceAll struct {
	DeviceId     string          `json:"deviceId"`
	BatteryLevel int             `json:"batteryLevel"`
	SignalStatus string          `json:"signalStatus"`
	IsLocked     bool            `json:"isLocked"`
	Status       string          `json:"status"`
	Location     *DeviceLocation `json:"location"`
}
