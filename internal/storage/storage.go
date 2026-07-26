package storage

import "pvz/internal/domain"

type Storage interface {
	GetOrder(orderID string) (domain.Order, error)
	DeleteOrder(orderID string) error
	SaveOrder(order domain.Order) error
}
