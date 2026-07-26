package domain

import "time"

type OrderStatus string

const (
	StatusAccepted OrderStatus = "ACCEPTED"
	StatusIsseud   OrderStatus = "ISSEUD"
	StatusReturned OrderStatus = "RETURNED"
)

type Order struct {
	OrderID    string      `json:"order_id"`
	CustomerID string      `json:"customer_id"`
	Status     OrderStatus `json:"status"`
	ExiredDate time.Time   `json:"expiry_date"`
	IssuedAt   time.Time   `json:"issued_at,omitempty"`
	UpdatedAt  time.Time   `json:"updated_at"`
}
