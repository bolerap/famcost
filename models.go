package main

type Cost struct {
	Id             int64   `json:"id"`
	ElectricAmount int64   `json:"electric_amount"`
	ElectricPrice  float64 `json:"electric_price"`
	WaterAmount    int64   `json:"water_amount"`
	WaterPrice     float64 `json:"water_price"`
	CheckedDate    string  `json:"checked_date"`
}

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     int64  `json:"role"`
}
