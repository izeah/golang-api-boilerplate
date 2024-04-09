package minioservice

import (
	"fmt"
	"path/filepath"
	"slices"

	"boilerplate/internal/abstraction"
	"boilerplate/pkg/mime"
	minioConstant "boilerplate/pkg/minio/constant"
	miniodto "boilerplate/pkg/minio/dto"
	"boilerplate/pkg/util/response"

	"github.com/google/uuid"
	goMinio "github.com/minio/minio-go/v7"
)

// UploadFile ...
func (s *service) UploadFile(ctx *abstraction.Context, payload *miniodto.MinioUploadFileRequest) (string, error) {
	file, err := ctx.FormFile("file")
	if err != nil {
		return "", response.CustomErrorBuilder(500, err.Error(), "upload_file")
	}

	if payload == nil || payload.Bucket == "" {
		return "", response.CustomErrorBuilder(400, "bucket is not set", "bucket is not set")
	}

	if exist := slices.Contains(s.Bucket, payload.Bucket); !exist {
		return "", response.CustomErrorBuilder(400, "bucket is not set", "bucket is not set")
	}

	fileSizeLimit := minioConstant.DefaultFileSize
	if payload.MaxFileSize > 0 {
		fileSizeLimit = payload.MaxFileSize
	}

	if file.Size > fileSizeLimit {
		return "", response.CustomErrorBuilder(400, "file size is too large", "file_size must be less than 1MB")
	}

	fileBody, err := file.Open()
	if err != nil {
		return "", response.CustomErrorBuilder(500, err.Error(), "open_file")
	}

	var isValidMimeType *bool
	if isValidMimeType, _, err = mime.AllowedMimeTypeImages(ctx, fileBody); err != nil || !*isValidMimeType {
		return "", response.CustomErrorBuilder(400, "file type is not allowed", "file_type")
	}

	// add this after read
	fileBody, err = file.Open()
	if err != nil {
		return "", response.CustomErrorBuilder(500, err.Error(), "open_file")
	}

	// Generate a unique file name
	extension := filepath.Ext(file.Filename)
	uniqueFileName := fmt.Sprintf("%s-%s%s", payload.Bucket, uuid.New().String(), extension)

	if _, err = s.Client.PutObject(
		ctx.Request().Context(),
		payload.Bucket,
		uniqueFileName,
		fileBody,
		file.Size,
		goMinio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		},
	); err != nil {
		return "", response.CustomErrorBuilder(500, err.Error(), "minio_upload")
	}

	return uniqueFileName, nil
}
