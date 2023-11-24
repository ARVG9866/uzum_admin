package models

import "time"

type Order struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`

	Coordinate *Coordinate `json:"coordinate"`

	Meta         string    `json:"meta"`
	Status       string    `json:"status"`
	DeliveryTime time.Time `json:"delivery_time"`
}

type Coordinate struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Statistic struct {
	CountSold int32      `json:"count_sold"`
	Earned    float64    `json:"earned"`
	Products  []*Product `json:"products"`
}

type GetAllSetting struct {
	Page int32 `json:"page"`
	Size int32 `json:"size"`
}
