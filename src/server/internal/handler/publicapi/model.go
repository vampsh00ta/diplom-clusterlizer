package publicapi

const (
	errNoFiles  = "NO_FILES"
	errReadForm = "CANT_READ_FORM"
	errReadFile = "CANT_READ_FILE"
)

type errResponse struct {
	Error string `json:"error"`
}
