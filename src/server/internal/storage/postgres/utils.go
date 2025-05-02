package postgresrep

import (
	"clusterlizer/pkg/utils"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

func patchQueryOptional[T any](
	query sq.UpdateBuilder,
	fieldName string,
	optValue utils.Optional[T],
) sq.UpdateBuilder {
	if optValue.Valid {
		query = query.Set(fieldName, optValue.Value)
	}
	return query
}

var (
	errNoRows = fmt.Errorf("no result")
)
