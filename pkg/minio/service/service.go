package minioservice

import (
	"boilerplate/pkg/minio"

	goMinio "github.com/minio/minio-go/v7"
)

type service struct {
	Client *goMinio.Client
	Bucket []string
}

func New(client *goMinio.Client, bucket []string) minio.Service {
	return &service{
		Client: client,
		Bucket: bucket,
	}
}
