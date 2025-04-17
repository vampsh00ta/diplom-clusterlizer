package document

import (
	"context"
)

type Service interface {
	SendToBroker(ctx context.Context, data any) error
}
