package postgresrep

import (
	"context"

	"clusterlizer/internal/entity"
	"clusterlizer/internal/storage/postgres/models"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

var cityFields = []string{
	fieldID,
	fieldName,
}

func (s *Storage) GetAllWithVacancyCount(ctx context.Context) ([]entity.CityWithVacancyCount, error) {
	client, err := s.db.Client(ctx)
	if err != nil {
		return nil, fmt.Errorf("client: %w", err)
	}

	//q := `
	//	SELECT s.name, COUNT(cj.vacancy_id) AS vacancy_count
	//	FROM city s
	//	LEFT JOIN vacancy_city cj ON s.id = cj.city_id
	//	GROUP BY s.id
	//	HAVING COUNT(cj.vacancy_id) > 0
	//	ORDER BY vacancy_count DESC;
	//`

	query := sq.Select(
		"s.name",
		"COUNT(cj.vacancy_id) AS vacancy_count",
	).
		From("city s").
		LeftJoin("vacancy_city cj ON s.id = cj.city_id").
		GroupBy("s.id").
		Having("COUNT(cj.vacancy_id) > 0").
		OrderBy("vacancy_count DESC")
	q, _, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query builder: %w", err)
	}
	rows, err := client.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("row query: %w", err)
	}
	defer rows.Close()

	rowModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.GetCitiesWithVacancyCount])
	if err != nil {
		return nil, err
	}

	res := make([]entity.CityWithVacancyCount, 0, len(rowModels))
	for _, rowModel := range rowModels {
		res = append(res, entity.CityWithVacancyCount{
			Name:         rowModel.Name,
			VacancyCount: rowModel.VacancyCount,
		})
	}

	return res, nil
}
