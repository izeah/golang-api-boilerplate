package minioservice

import (
	"context"
	"time"

	miniodto "boilerplate/pkg/minio/dto"
	"boilerplate/pkg/util/response"
)

// ObjectURL ...
func (s *service) ObjectURL(ctx context.Context, payload *miniodto.MinioObjectURLRequest) (string, error) {
	if payload == nil {
		return "", response.CustomErrorBuilder(400, "need filter", "need filter")
	}
	objectURL, err := s.Client.PresignedGetObject(ctx, payload.Bucket, payload.FileName, time.Minute*15, nil)
	if err != nil {
		return "", response.CustomErrorBuilder(500, err.Error(), "minio_get_object_url")
	}
	return objectURL.String(), nil
}
