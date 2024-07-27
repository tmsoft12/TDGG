package models

import "time"

type Driver struct {
	ID         int       `json:"id"`
	CreateTime time.Time `json:"create_time"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	CarNumber  string    `json:"car_number"`
	CarModel   string    `json:"car_model"`
	Weight     int       `json:"weight"`
	Country    string    `json:"country"`
}
