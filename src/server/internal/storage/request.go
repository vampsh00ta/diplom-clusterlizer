package storage

import (
	"clusterlizer/internal/entity"
	"clusterlizer/pkg/utils"
)

type CreateRequestParams struct {
	ID entity.RequestID `db:"id"`
}

type UpdateRequestParams struct {
	ID entity.RequestID `db:"id"`

	Result utils.Optional[*[]byte]       `db:"result"`
	Status utils.Optional[entity.Status] `db:"status"`
}
