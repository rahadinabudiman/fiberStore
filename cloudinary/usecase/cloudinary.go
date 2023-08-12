package usecase

import (
	"fiberStore/helpers"
	"fiberStore/models"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
)

type media struct{}

func NewMediaUpload() models.CloudinaryUsecase {
	return &media{}
}

func (*media) FileUpload(file models.File) (string, error) {
	//validate
	err := validate.Struct(file)
	if err != nil {
		return "", err
	}

	//upload
	uploadUrl, err := helpers.ImageUpload(file.File)
	if err != nil {
		return "", err
	}
	return uploadUrl, nil
}

func (*media) RemoteUpload(url models.Url) (string, error) {
	//validate
	err := validate.Struct(url)
	if err != nil {
		return "", err
	}

	//upload
	uploadUrl, errUrl := helpers.ImageUpload(url.Url)
	if errUrl != nil {
		return "", err
	}
	return uploadUrl, nil
}
