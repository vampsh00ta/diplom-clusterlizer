package entity

type FileID string

func (f FileID) String() string {
	return string(f)
}

type FileType int

const (
	FileTypePDF FileType = iota
	FileTypeDOCX
	FileTypeUnspecified
)

const (
	textFileTypePDF       = "pdf"
	textFileTypeDOCX      = "docx"
	UNSPECIFIED_FILE_TYPE = "unspecifed"
)

var AllowedTypes = []string{
	textFileTypePDF,
	textFileTypeDOCX,
}

func (s FileType) String() string {
	switch s {
	case FileTypePDF:
		return textFileTypePDF
	case FileTypeDOCX:
		return textFileTypeDOCX

	}
	return UNSPECIFIED_FILE_TYPE
}

func FileTypeFromString(s string) FileType {
	switch s {
	case textFileTypePDF:
		return FileTypePDF
	case textFileTypeDOCX:
		return FileTypeDOCX

	}
	return FileTypeUnspecified
}

type File struct {
	ID    FileID   `json:"id"`
	Key   string   `json:"key"`
	Type  FileType `json:"type"`
	Title string   `json:"title"`
}
