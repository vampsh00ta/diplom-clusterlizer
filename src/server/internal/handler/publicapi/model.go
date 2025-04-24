package publicapi

const (
	errNoFiles                 = "NO_FILES"
	errSeveralGroupCountForms  = "SEVERAL_GROUP_COUNT_FORMS"
	errLessFilesThanGroupCount = "LESS_FILES_THAN_GROUP_COUNT"

	errReadForm             = "CANT_READ_FORM"
	errReadFile             = "CANT_READ_FILE"
	errNowAllowedFileFormat = "NOW_ALLOWED_FILE_FORMAT"
)

type errResponse struct {
	Error string `json:"error"`
}
