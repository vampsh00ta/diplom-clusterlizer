package publicapi

import (
	"clusterlizer/internal/entity"
	"strings"
)

func correctFileFormat(file string) bool {
	splitedFile := strings.Split(file, ".")
	if len(splitedFile) < 2 {
		return false
	}
	format := splitedFile[len(splitedFile)-1]
	for _, allowedFormat := range entity.AllowedTypes {
		if allowedFormat == format {
			return true
		}
	}
	return false
}
