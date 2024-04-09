package config

import (
	"os"
	"strings"
	"sync"

	"boilerplate/pkg/circuitbreaker"
	"boilerplate/pkg/util/priority"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioConfig struct {
	Host        string
	AccessKey   string
	SecretKey   string
	Bucket      []string
	MinioClient *minio.Client
}

var (
	minioConfig *MinioConfig
	minioOnce   sync.Once
)

func Minio() *MinioConfig {
	minioOnce.Do(func() {
		minioConfig = new(MinioConfig)

		minioConfig.Host = priority.PriorityString(os.Getenv("MINIO_HOST"))
		minioConfig.SecretKey = priority.PriorityString(os.Getenv("MINIO_SECRET_KEY"))
		minioConfig.AccessKey = priority.PriorityString(os.Getenv("MINIO_ACCESS_KEY"))
		minioConfig.Bucket = priority.PrioritySliceString(strings.Split(os.Getenv("MINIO_BUCKET"), ","))

		var err error
		if minioConfig.MinioClient, err = minio.New(minioConfig.Host, &minio.Options{
			Creds:     credentials.NewStaticV4(minioConfig.AccessKey, minioConfig.SecretKey, ""),
			Secure:    true,
			Transport: circuitbreaker.NewClient().Transport,
		}); err != nil {
			panic(err)
		}
	})
	return minioConfig
}
