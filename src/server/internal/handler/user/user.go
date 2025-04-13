package user

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
}

func New() *Handler {
	return &Handler{}
}

type saveFilterRequest struct {
	UserTgID int `json:"tg_id"`
	// JobName    string   `json:"job_name"`
	City       string `json:"city"`
	Experience string `json:"experience"`
	Company    string `json:"company"`
	Speciality string `json:"speciality"`

	KeywordsSlugs []string `json:"keywords"`
}
type saveFilterResponse struct {
	Status string `json:"status"`
}

func (u Handler) SaveFilter(ctx *fiber.Ctx) error {
	var req saveFilterRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	res := saveFilterResponse{"ok"}
	return ctx.Status(fiber.StatusCreated).JSON(res)
}
