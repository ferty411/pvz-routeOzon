package service

import (
	"errors"
	"pvz/internal/domain"
	"pvz/internal/storage"
	"time"
)

type PVZService struct {
	storage storage.Storage
}

func NewPvzService(storage storage.Storage) *PVZService {
	return &PVZService{
		storage: storage,
	}
}

func (p *PVZService) AcceptOrder(orderID, customerID string, expireDate time.Time) error {
	if expireDate.Before(time.Now()) {
		return domain.ErrStorageExpired
	}

	_, err := p.storage.GetOrder(orderID)
	if err == nil {
		return domain.ErrOrderAlreadyExists
	}

	if !errors.Is(err, domain.ErrOrderNotFound) {
		return err
	}

	NewOrder := domain.Order{
		OrderID:    orderID,
		CustomerID: customerID,
		Status:     domain.StatusAccepted,
		ExiredDate: expireDate,
		UpdatedAt:  time.Now(),
	}

	return p.storage.SaveOrder(NewOrder)
}
