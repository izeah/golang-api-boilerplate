package factory

import (
	"boilerplate/internal/config"
	"boilerplate/internal/repository"
	"boilerplate/pkg/database"
	"boilerplate/pkg/redis"

	"github.com/minio/minio-go/v7"
	goRedis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Factory struct {
	MinioClient *minio.Client
	RedisClient *goRedis.Client

	DB             *gorm.DB
	UserRepository repository.User
}

func NewFactory() *Factory {
	f := new(Factory)
	f.SetupDB()
	f.SetupClient()
	f.SetupRepository()
	return f
}

func (f *Factory) SetupDB() {
	f.DB = database.PSQL()
}

func (f *Factory) SetupClient() {
	f.RedisClient = redis.Client()
	f.MinioClient = config.Minio().MinioClient
}

func (f *Factory) SetupRepository() {
	if f.DB == nil {
		panic("Failed setup repository, db is undefined")
	}

	f.UserRepository = repository.NewUser(f.DB)
}
