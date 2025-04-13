package models

type GetAll struct {
	ID             int    `json:"id" db:"id"`
	SlugID         string `json:"slug_id" db:"slug_id"`
	Name           string `json:"name" db:"name"`
	City           string `json:"city_name" db:"city_name"`
	Category       string `json:"category" db:"category"`
	Description    string `json:"description" db:"description"`
	CompanySlug    string `json:"company_slug" db:"company_slug"`
	CompanyName    string `json:"company_name" db:"company_name"`
	Link           string `json:"link" db:"link"`
	ExperienceSlug string `json:"experience_slug" db:"experience_slug"`
	ExperienceName string `json:"experience_name" db:"experience_name"`
	ExperienceID   string `json:"experience_id" db:"experience_id"`
}
type GetAllWithFilter struct {
	ID          int    `json:"id" db:"id"`
	SlugID      string `json:"slug_id" db:"slug_id"`
	Name        string `json:"name" db:"name"`
	City        string `json:"city_name" db:"city_name"`
	Category    string `json:"category" db:"category"`
	Description string `json:"description" db:"description"`
	CompanySlug string `json:"company_slug" db:"company_slug"`
	CompanyName string `json:"company_name" db:"company_name"`
	Link        string `json:"link" db:"link"`

	ExperienceSlug string `json:"experience_slug" db:"experience_slug"`
	ExperienceName string `json:"experience_name" db:"experience_name"`

	SpecialitySlug *string `json:"speciality_slug" db:"speciality_slug"`
	SpecialityName *string `json:"speciality_name" db:"speciality_name"`

	KeywordSlug *string `json:"keyword_slug" db:"keyword_slug"`
	KeywordName *string `json:"keyword_name" db:"keyword_name"`

	ExperienceID string  `json:"experience_id" db:"experience_id"`
	Rank         float64 `json:"rank" db:"rank"`
}
