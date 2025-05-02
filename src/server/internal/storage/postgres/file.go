package storage

import "time"

type File struct {
	ID        FileID
	Result    GraphData
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}
