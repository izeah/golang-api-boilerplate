package config

import (
	"os"
	"sync"
	"time"

	"boilerplate/pkg/util/priority"
)

type JwtConfig struct {
	JWTKey             string // for access_token
	JWTRefKey          string // for refresh_token
	EncryptionKey      string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

var (
	jwt     *JwtConfig
	jwtOnce sync.Once
)

func JWT() *JwtConfig {
	jwtOnce.Do(func() {
		jwt = &JwtConfig{
			JWTKey:             priority.PriorityString(os.Getenv("JWT_KEY")),
			JWTRefKey:          priority.PriorityString(os.Getenv("JWT_REF_KEY")),
			EncryptionKey:      priority.PriorityString(os.Getenv("ENC_KEY")),
			AccessTokenExpiry:  time.Duration(5 * time.Minute),
			RefreshTokenExpiry: time.Duration(15 * time.Minute),
		}
	})
	return jwt
}
