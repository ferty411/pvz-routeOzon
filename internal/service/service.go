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

func (p *PVZService) ReturnCorier(orderID string) error {
	order, err := p.storage.GetOrder(orderID)
	if err != nil {
		return domain.ErrOrderNotFound
	}

	if order.Status != domain.StatusAccepted {
		return domain.ErrInvalidStatus
	}

	if order.ExiredDate.After(time.Now()) {
		return domain.ErrValidationFailed
	}
	return p.storage.DeleteOrder(orderID)
}

func (p *PVZService) IssueOrders(orderIDs []string, customerID string) error {
	now := time.Now()
	var ordersToSave []domain.Order
	for _, ord := range orderIDs {
		order, err := p.storage.GetOrder(ord)
		if err != nil {
			return domain.ErrOrderNotFound
		}
		if order.CustomerID != customerID {
			return domain.ErrValidationFailed
		}
		if order.Status != domain.StatusAccepted {
			return domain.ErrInvalidStatus
		}
		if order.ExiredDate.Before(now) {
			return domain.ErrStorageExpired
		}

		order.Status = domain.StatusIsseud
		order.IssuedAt = now
		order.UpdatedAt = now
		ordersToSave = append(ordersToSave, order)
	}

	for _, val := range ordersToSave {
		if err := p.storage.SaveOrder(val); err != nil {
			return err
		}
	}
	return nil
}

func (p *PVZService) ProccessClientAction(action string, customerID string, ordersID []string) error {
	if len(ordersID) == 0 {
		return domain.ErrValidationFailed
	}

	switch action {
	case "issue":
		return p.IssueOrders(ordersID, customerID)
	case "return":
		return p.ReturnFromClient(ordersID, customerID)

	default:
		return errors.New("неизвестное действие: используйте 'issue' или 'return'")
	}
}
