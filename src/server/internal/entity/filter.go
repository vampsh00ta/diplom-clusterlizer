package entity

type Filter struct {
	// ID         int      `json:"id"`
	UserTgID int    `json:"tg_id"`
	JobName  string `json:"job_name"`
	City     string `json:"city"`
}

type AllFilter struct {
	Cities []CityWithVacancyCount `json:"cities"`
}
