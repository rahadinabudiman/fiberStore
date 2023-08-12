package helpers

import (
	"context"
	"errors"
	"fiberStore/author"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func ImageUpload(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cloud, err := cloudinary.NewFromParams(author.EnvCloudName(), author.EnvCloudAPIKey(), author.EnvCloudAPISecret())
	if err != nil {
		return "", errors.New("failed to connect instance")
	}

	// Upload Images
	uploadParams, err := cloud.Upload.Upload(ctx, input, uploader.UploadParams{Folder: author.EnvCloudUploadFolder()})
	if err != nil {
		return "", errors.New("failed to upload image")
	}

	return uploadParams.SecureURL, nil
}
