package common

import "errors"

var (
	ErrInvalidOrderNumber     = errors.New("invalid order number")
	ErrOrderAlreadyAdded      = errors.New("order already added by user")
	ErrOrderBelongAnotherUser = errors.New("order belong another user")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrInvalidCredentials     = errors.New("invalid credentials")
	ErrPaymentInsufficient    = errors.New("payment insufficient")
	ErrNoContent              = errors.New("no content")
)
