package models

type GetCitiesWithVacancyCount struct {
	Name         string `json:"name" db:"name"`
	VacancyCount int    `json:"vacancy_count" db:"vacancy_count"`
}
