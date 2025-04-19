package publicapi

import (
	"clusterlizer/internal/entity"
	documentsrvc "clusterlizer/internal/service/document"
	requestsrvc "clusterlizer/internal/service/request"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	log          *zap.SugaredLogger
	documentSrvc documentsrvc.Service
	requestSrvc  requestsrvc.Service
}

func New(
	log *zap.SugaredLogger,
	documentSrvc documentsrvc.Service,
	requestSrvc requestsrvc.Service,

) *Handler {
	return &Handler{
		log:          log,
		documentSrvc: documentSrvc,
		requestSrvc:  requestSrvc,
	}
}

type uploadFilesRequest struct {
}
type uploadFilesResponse struct {
	UUID uuid.UUID `json:"uuid"`
}

func (h *Handler) UploadFiles(ctx *fiber.Ctx) error {
	var req uploadFilesRequest
	_ = req
	uuid := uuid.New()
	// res := uploadFilesResponse{UUID: uuid.New()}
	if err := h.requestSrvc.CreateRequest(ctx.Context(), requestsrvc.CreateRequestParams{
		ID: entity.RequestID(uuid.String()),
	}); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errResponse{
			Error: err.Error(),
		})
	}
	res:=uploadFilesResponse{UUID: uuid}
	return ctx.Status(fiber.StatusCreated).JSON(res)
}

type getClusterizationsRequest struct {
	UUID uuid.UUID `json:"uuid"`
}
type getClusterizationsResponse struct {
	UUID uuid.UUID `json:"uuid"`
}

func (h *Handler) GetClusterizations(ctx *fiber.Ctx) error {
	var req getClusterizationsRequest
	_ = req
	if err := ctx.QueryParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errResponse{
			Error: err.Error(),
		})
	}
	_ = req

	res := getClusterizationsResponse{UUID: req.UUID}
	return ctx.Status(fiber.StatusCreated).JSON(res)
}

type GetCurrentQueueRequest struct {
	UUID uuid.UUID `json:"uuid"`
}
type GetCurrentQueueResponse struct {
	Number int `json:"number"`
}

func (h *Handler) GetCurrentQueue(ctx *fiber.Ctx) error {
	var req getClusterizationsRequest
	_ = req
	if err := ctx.QueryParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(errResponse{
			Error: err.Error(),
		})
	}
	_ = req

	res := getClusterizationsResponse{}
	return ctx.Status(fiber.StatusCreated).JSON(res)
}
