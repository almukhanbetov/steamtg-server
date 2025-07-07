package models

type Driver struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	IIN    string  `json:"iin"`
	Photo  string  `json:"photo"`
	Lon    float64 `json:"lon"`
	Lat    float64 `json:"lat"`
	CarID  int     `json:"car_id"`
	Phone  string  `json:"phone,omitempty"`
}

type Category struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Client struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Phone       string  `json:"phone"`
	Location    string  `json:"location"`
	PlateNumber string  `json:"plate_number"`
	CategoryID  int     `json:"category_id"`
	Lon         float64 `json:"lon,omitempty"`
	Lat         float64 `json:"lat,omitempty"`
}

type Order struct {
	ID          int     `json:"id"`
	ClientID    int     `json:"client_id"`
	DriverID    int     `json:"driver_id"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	ClientName  string  `json:"client_name,omitempty"`
	ClientPhone string  `json:"client_phone,omitempty"`
	Lon         float64 `json:"lon,omitempty"`
	Lat         float64 `json:"lat,omitempty"`
}
