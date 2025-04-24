package document

import (
	"context"
)

type Service interface {
	SendDocumentNames(ctx context.Context, params SendDocumentParams) error
}

type SendDocumentParams struct {
	GroupCount int      `json:"group_count"`
	Keys       []string `json:"names"`
}
