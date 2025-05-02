package entity

import "fmt"

var (
	ErrNoResult  = fmt.Errorf("no result")
	ErrReqFailed = fmt.Errorf("request finished with an error")
)
