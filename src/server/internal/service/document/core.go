package document

import (
	"context"
)

type Service interface {
	SendDocumentNames(ctx context.Context, names []string) error
}
