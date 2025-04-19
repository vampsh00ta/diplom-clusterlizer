package entity

import "time"

type RequestID string

type Request struct {
	ID        RequestID
	Result    []byte
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}
