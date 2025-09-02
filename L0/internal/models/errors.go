package models

import "errors"

// Ошибки валидации
var (
	ErrInvalidOrderUID    = errors.New("invalid order UID")
	ErrInvalidTrackNumber = errors.New("invalid track number")
	ErrNoItems            = errors.New("order must contain at least one item")
	ErrOrderNotFound      = errors.New("order not found")
	ErrDatabaseConnection = errors.New("database connection error")
	ErrKafkaConnection    = errors.New("kafka connection error")
)
