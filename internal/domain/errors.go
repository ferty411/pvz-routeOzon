package domain

import "errors"

var (
	ErrOrderAlreadyExists = errors.New("ORDER_ALREADY_EXISTS")
	ErrOrderNotFound      = errors.New("ORDER_NOT_FOUND")
	ErrStorageExpired     = errors.New("STORAGE_EXPIRED")
	ErrValidationFailed   = errors.New("VALIDATION_FAILED")

	ErrInvalidStatus              = errors.New("INVALID_STATUS")
	ErrReturnPeriodExpired        = errors.New("RETURN_PERIOD_EXPIRED")
	ErrOrdersFromDifferentClients = errors.New("ORDERS_FROM_DIFFERENT_CLIENTS")
	ErrWeightTooHeavy             = errors.New("WEIGHT_TOO_HEAVY")
)
