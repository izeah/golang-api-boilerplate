package minio

import (
	"context"

	"boilerplate/internal/abstraction"
	miniodto "boilerplate/pkg/minio/dto"
)

type Service interface {
	UploadFile(ctx *abstraction.Context, payload *miniodto.MinioUploadFileRequest) (string, error)
	ObjectURL(ctx context.Context, payload *miniodto.MinioObjectURLRequest) (string, error)
}
