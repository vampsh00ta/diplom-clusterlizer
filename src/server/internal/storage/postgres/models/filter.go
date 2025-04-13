package models

type Filter struct {
	ID       int `json:"id" db:"id"`
	UserTgID int `json:"tg_id" db:"tg_id"`
	// JobName    string   `json:"job_name"`
	City           *string `json:"city" db:"city"`
	ExperienceSlug *string `json:"experience_slug" db:"experience_slug"`
	ExperienceName *string `json:"experience_name" db:"experience_name"`

	CompanySlug *string `json:"company_slug" db:"company_slug"`
	CompanyName *string `json:"company_name" db:"company_name"`

	KeywordSlug *string `json:"keyword_slug" db:"keyword_slug"`
	KeywordName *string `json:"keyword_name" db:"keyword_name"`

	SpecialitySlug *string `json:"speciality_slug" db:"speciality_slug"`
	SpecialityName *string `json:"speciality_name" db:"speciality_name"`
}
