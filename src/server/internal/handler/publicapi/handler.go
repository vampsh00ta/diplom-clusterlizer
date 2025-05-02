package publicapi

import (
	"bytes"
	"clusterlizer/internal/entity"
	documentsrvc "clusterlizer/internal/service/document"
	filesrvc "clusterlizer/internal/service/file"
	requestsrvc "clusterlizer/internal/service/request"

	s3 "clusterlizer/internal/service/s3"
	"strconv"
	"time"

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
	fileSrvc     filesrvc.Service
}

func New(
	log *zap.SugaredLogger,
	documentSrvc documentsrvc.Service,
	requestSrvc requestsrvc.Service,
	s3Srvc s3.Service,
	fileSrvc filesrvc.Service,

) *Handler {
	return &Handler{
		log:          log,
		documentSrvc: documentSrvc,
		requestSrvc:  requestSrvc,
		s3Srvc:       s3Srvc,
		fileSrvc:     fileSrvc,
	}
}

type uploadFilesRequest struct {
	GroupCount int
	Files      [][]byte
}

type uploadFilesResponse struct {
	UUID uuid.UUID `json:"uuid"`
}

// похорошему нужно добавить temporal, чтобы добиться какой-никакой транзитивности
func (h *Handler) UploadFiles(ctx *fiber.Ctx) error {
	id := uuid.New()

	req, err := h.getUploadFilesFromForm(ctx)
	if err != nil {
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
	fileNames := make([]string, 0, len(req.Files))
	for i, file := range req.Files {
		fileKey := fmt.Sprintf("%s_%s", id.String(), strconv.Itoa(i))
		fileNames = append(fileNames, fileKey)
		if err := h.s3Srvc.Upload(ctx.Context(), fileKey, file); err != nil {
			h.log.Error(err)

			return ctx.Status(fiber.StatusInternalServerError).JSON(errResponse{
				Error: err.Error(),
			})
		}
	}

	if err := h.documentSrvc.SendDocumentNames(ctx.Context(), documentsrvc.SendDocumentParams{
		GroupCount: req.GroupCount,
		Keys:       fileNames,
	}); err != nil {
		h.log.Error(err)

		return ctx.Status(fiber.StatusInternalServerError).JSON(errResponse{
			Error: err.Error(),
		})
	}
	res := uploadFilesResponse{UUID: id}
	return ctx.Status(fiber.StatusCreated).JSON(res)
}

type getClusterizationsResponse struct {
	ID        string           `json:"id"`
	Result    entity.GraphData `json:"result"`
	Status    string           `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func (h *Handler) GetClusterizations(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")

	ID, err := uuid.Parse(idParam)
	if err != nil {
		h.log.Error(err)

		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid UUID",
		})
	}
	res, err := h.requestSrvc.GetRequestByIDDone(ctx.Context(), entity.RequestID(ID.String()))
	if err != nil {
		h.log.Error(err)
		err = handleError(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(errResponse{
			Error: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(getClusterizationsResponse{
		ID:        res.ID.String(),
		Result:    res.Result,
		Status:    res.Status.String(),
		UpdatedAt: res.UpdatedAt,
		CreatedAt: res.CreatedAt,
	})
}

type downloadFileRequest struct {
	Key uuid.UUID `json:"key"`
}

func (h *Handler) DownloadFile(ctx *fiber.Ctx) error {
	key := ctx.Params("key")

	file, err := h.fileSrvc.GetRequestByKey(ctx.Context(), key)
	if err != nil {
		h.log.Error(err)

		err = handleError(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(errResponse{
			Error: err.Error(),
		})
	}
	fullFileName := fmt.Sprintf("%s.%s", file.Title, file.Type.String())
	fileBytes, err := h.s3Srvc.Download(ctx.Context(), key)
	if err != nil {
		h.log.Error(err)

		err = handleError(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(errResponse{
			Error: err.Error(),
		})
	}
	ctx.Attachment(fullFileName)
	return ctx.Send(fileBytes)
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

func (h *Handler) getUploadFilesFromForm(ctx *fiber.Ctx) (uploadFilesRequest, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		h.log.Error(err)

		return uploadFilesRequest{}, fmt.Errorf(errReadForm)
	}

	filesForm := form.File["file"]
	if len(filesForm) == 0 {
		h.log.Error(fmt.Errorf(errNoFiles))

		return uploadFilesRequest{}, fmt.Errorf(errNoFiles)
	}
	files := make([][]byte, 0, len(filesForm))
	for _, fileForm := range filesForm {
		fileBytes, err := fileFormToBytes(fileForm)
		if err != nil {
			h.log.Error(fmt.Errorf("file form to bytes: %w", err))

			return uploadFilesRequest{}, fmt.Errorf("file form to bytes: %w", err)
		}
		if !correctFileFormat(fileForm.Filename) {
			h.log.Error(fmt.Errorf("%s; file name:%s", errNowAllowedFileFormat, fileForm.Filename))

			return uploadFilesRequest{}, fmt.Errorf(errNowAllowedFileFormat)
		}
		files = append(files, fileBytes)
	}
	groupCountForm := form.Value["group_count"]
	if len(groupCountForm) != 1 {
		h.log.Error(fmt.Errorf(errSeveralGroupCountForms))

		return uploadFilesRequest{}, fmt.Errorf(errSeveralGroupCountForms)
	}
	groupCountStr := groupCountForm[0]
	groupCount, err := strconv.Atoi(groupCountStr)
	if err != nil {
		h.log.Error(err)

		return uploadFilesRequest{}, err
	}

	if len(files) < groupCount {
		h.log.Error(fmt.Errorf(errLessFilesThanGroupCount))

		return uploadFilesRequest{}, fmt.Errorf(errLessFilesThanGroupCount)
	}

	return uploadFilesRequest{
		Files:      files,
		GroupCount: groupCount,
	}, nil
}
