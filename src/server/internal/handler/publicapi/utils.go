package publicapi

import "strings"

var (
	allowedFileFormats = []string{
		"pdf",
		"docx",
	}
)

func correctFileFormat(file string) bool {
	splitedFile := strings.Split(file, ".")
	if len(splitedFile) < 2 {
		return false
	}
	format := splitedFile[len(splitedFile)-1]
	for _, allowedFormat := range allowedFileFormats {
		if allowedFormat == format {
			return true
		}
	}
	return false
}
