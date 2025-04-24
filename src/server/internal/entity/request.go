package entity

import "time"

type RequestID string

func (r RequestID) String() string {
	return string(r)
}

type Groups []Group

type Request struct {
	ID        RequestID
	Result    Groups
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}
