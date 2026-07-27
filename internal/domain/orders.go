package domain

import (
	"time"
)

type OrderStatus string
type PackagingType string

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
	Price      float64     `json:"price"`
	Weight     float64     `json:"weight"`
	Packaging  []string    `json:"packaging"`
}

const (
	PackBox  PackagingType = "box"
	PackBag  PackagingType = "bag"
	PackFilm PackagingType = "film"
)

const (
	BagPrice       = 5.0
	BagWeightLimit = 10.0

	BoxPrice       = 20.0
	BoxWeightLimit = 30.0

	FilmPrice = 1.0
)

type OrderPackeger interface {
	GetPrice() float64
	ValidateWeight(weight float64) error
}

type BaseOrder struct {
	InitialPrice float64
}

func (b *BaseOrder) GetPrice() float64                   { return b.InitialPrice }
func (b *BaseOrder) ValidateWeight(weight float64) error { return nil }

type BagPackaging struct {
	Wrapper OrderPackeger
}

func (b *BagPackaging) GetPrice() float64 {
	return b.GetPrice() + BagPrice
}

func (b *BagPackaging) ValidateWeight(weight float64) error {
	if weight >= BagWeightLimit {
		return ErrWeightTooHeavy
	}
	return b.Wrapper.ValidateWeight(weight)
}

type BoxPackaging struct {
	Wrapper OrderPackeger
}

func (b *BoxPackaging) GetPrice() float64 {
	return b.GetPrice() + BoxPrice
}

func (b *BoxPackaging) ValidateWeight(weight float64) error {
	if weight >= BoxWeightLimit {
		return ErrWeightTooHeavy
	}
	return b.Wrapper.ValidateWeight(weight)
}

type FilmPackaging struct {
	Wrapper OrderPackeger
}

func (f *FilmPackaging) GetPrice() float64 {
	return f.Wrapper.GetPrice() + FilmPrice
}

func (f *FilmPackaging) ValidateWeight(weight float64) error {
	return f.Wrapper.ValidateWeight(weight)
}
