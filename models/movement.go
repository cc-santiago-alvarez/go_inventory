package models

import "time"

type MovementType string

const (
	MovementEntry MovementType = "entry"
	MovementExit  MovementType = "exit"
)

type Movement struct {
	ID        string          `json:"id"`
	Product   ProductCategory `json:"product"`
	Type      MovementType    `json:"type"`
	Quantity  int             `json:"quantity"`
	Reason    string          `json:"reason"`
	Date      time.Time       `json:"date"`
	CreatedAt time.Time       `json:"created_at"`
}

type CreateMovementRequest struct {
	ProductID string       `json:"product_id"`
	Type      MovementType `json:"type"`
	Quantity  int          `json:"quantity"`
	Reason    string       `json:"reason"`
}
