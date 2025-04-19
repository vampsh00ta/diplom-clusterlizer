package entity

type Status int

const (
	StatusCreated Status = iota
	StatusProcessing
	StatusDone
	StatusError
	StatusUnspecified
)

const (
	textStatusCreated    = "created"
	textStatusProcessing = "processing"
	textStatusDone       = "done"
	textStatusError      = "error"
	UNSPECIFIED_STATUS   = "UNSPECIFIED_STATUS"
)

func (s Status) String() string {
	switch s {
	case StatusCreated:
		return textStatusCreated
	case StatusProcessing:
		return textStatusProcessing
	case StatusDone:
		return textStatusDone
	case StatusError:
		return UNSPECIFIED_STATUS
	}
	return UNSPECIFIED_STATUS
}

func StatusFromString(s string) Status {
	switch s {
	case textStatusCreated:
		return StatusCreated
	case textStatusProcessing:
		return StatusProcessing
	case textStatusDone:
		return StatusDone
	case textStatusError:
		return StatusError
	}
	return StatusUnspecified
}
