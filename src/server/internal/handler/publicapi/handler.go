package publicapi

import (
	"bytes"
	"clusterlizer/internal/entity"
	documentsrvc "clusterlizer/internal/service/document"
	requestsrvc "clusterlizer/internal/service/request"
	s3 "clusterlizer/internal/service/s3"
	"strconv"

	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"mime/multipart"
)

type Handler struct {
	log          *zap.SugaredLogger
	documentSrvc documentsrvc.Service
	requestSrvc  requestsrvc.Service
	s3Srvc       s3.Service
}

func New(
	log *zap.SugaredLogger,
	documentSrvc documentsrvc.Service,
	requestSrvc requestsrvc.Service,
	s3Srvc s3.Service,
) *Handler {
	return &Handler{
		log:          log,
		documentSrvc: documentSrvc,
		requestSrvc:  requestSrvc,
		s3Srvc:       s3Srvc,
	}
}

type uploadFilesResponse struct {
	UUID uuid.UUID `json:"uuid"`
}

// похорошему нужно добавить temporal, чтобы добиться какой-никакой транзитивности
func (h *Handler) UploadFiles(ctx *fiber.Ctx) error {
	id := uuid.New()

	files, err := h.getFiles(ctx)
	if err != nil {
		h.log.Error(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(errResponse{
			Error: err.Error(),
		})
	}

	if err := h.requestSrvc.CreateRequest(ctx.Context(), requestsrvc.CreateRequestParams{
		ID: entity.RequestID(id.String()),
	}); err != nil {
		h.log.Error(err)

		return ctx.Status(fiber.StatusInternalServerError).JSON(errResponse{
			Error: err.Error(),
		})
	}
	fileNames := make([]string, 0, len(files))
	for i, file := range files {
		fileKey := fmt.Sprintf("%s_%s", id.String(), strconv.Itoa(i))
		fileNames = append(fileNames, fileKey)
		if err := h.s3Srvc.Upload(ctx.Context(), fileKey, file); err != nil {
			h.log.Error(err)

			return ctx.Status(fiber.StatusInternalServerError).JSON(errResponse{
				Error: err.Error(),
			})
		}
	}

	if err := h.documentSrvc.SendDocumentNames(ctx.Context(), fileNames); err != nil {
		h.log.Error(err)

		return ctx.Status(fiber.StatusInternalServerError).JSON(errResponse{
			Error: err.Error(),
		})
	}
	res := uploadFilesResponse{UUID: id}
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

func fileFormToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return nil, fmt.Errorf("read from: %w", err)
	}
	return buf.Bytes(), nil
}

func (h *Handler) getFiles(ctx *fiber.Ctx) ([][]byte, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		h.log.Error(err)
		return nil, fmt.Errorf(errReadForm)
	}

	filesForm := form.File["file"]
	if len(filesForm) == 0 {
		return nil, fmt.Errorf(errNoFiles)
	}
	files := make([][]byte, 0, len(filesForm))
	for _, fileForm := range filesForm {
		fileBytes, err := fileFormToBytes(fileForm)
		if err != nil {
			return nil, fmt.Errorf("file form to bytes: %w", err)
		}
		files = append(files, fileBytes)
	}
	return files, nil
}
