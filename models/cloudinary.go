package models

import "mime/multipart"

type File struct {
	File multipart.File `json:"file,omitempty" validate:"required"`
}

type Url struct {
	Url string `json:"url,omitempty" validate:"required"`
}

type CloudinaryUsecase interface {
	FileUpload(file File) (string, error)
	RemoteUpload(url Url) (string, error)
}
