package document

import (
	"context"
)

type Service interface {
	SendDocumentNames(ctx context.Context, params SendDocumentParams) error
}

type SendDocumentParams struct {
	Keys []string `json:"keys"`
}
