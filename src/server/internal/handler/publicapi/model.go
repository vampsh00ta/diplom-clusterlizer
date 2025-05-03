package publicapi

import (
	"errors"
	"fmt"

	"clusterlizer/internal/entity"
)

const (
	errNoFiles                 = "NO_FILES"
	errSeveralGroupCountForms  = "SEVERAL_GROUP_COUNT_FORMS"
	errLessFilesThanGroupCount = "LESS_FILES_THAN_GROUP_COUNT"
	errNoResult                = "NO_RESULT"
	errReadForm                = "CANT_READ_FORM"
	errReadFile                = "CANT_READ_FILE"
	errNowAllowedFileFormat    = "NOW_ALLOWED_FILE_FORMAT"
	errRequestFailed           = "REQUEST_FAILED"
)

type errResponse struct {
	Error string `json:"error"`
}

func handleError(err error) error {
	switch {
	case errors.Is(err, entity.ErrNoResult):
		return fmt.Errorf(errNoResult)
	case errors.Is(err, entity.ErrReqFailed):
		return fmt.Errorf(errRequestFailed)

	}
	return err
}
