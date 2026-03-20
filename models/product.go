package models

import "time"

type ProductCategory struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Prefix string `json:"prefix"`
	Description string          `json:"description"`
}

type Product struct {
	ID          string      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Category    ProductCategory `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
