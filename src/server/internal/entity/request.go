package entity

import "time"

type RequestID string

func (r RequestID) String() string {
	return string(r)
}

type Request struct {
	ID        RequestID
	Result    GraphData
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}
