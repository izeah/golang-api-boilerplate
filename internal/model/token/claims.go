package modeltoken

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"boilerplate/internal/abstraction"
	"boilerplate/internal/config"
	"boilerplate/pkg/util/aescrypt"

	"github.com/golang-jwt/jwt/v4"
)

type AccessTokenClaims struct {
	ID     string `json:"id"`
	RoleID string `json:"rid"`
	Exp    int64  `json:"exp"`

	jwt.RegisteredClaims
}

func (c AccessTokenClaims) AuthContext() (*abstraction.AuthContext, error) {
	var (
		id  int
		rid int
		err error

		encryptionKey = config.JWT().EncryptionKey
	)

	destructID := c.ID
	if destructID == "" {
		return nil, errors.New("invalid_access_token")
	}
	if id, err = strconv.Atoi(fmt.Sprintf("%v", destructID)); err != nil {
		if destructID, err = aescrypt.DecryptAES(fmt.Sprintf("%v", destructID), encryptionKey); err != nil {
			return nil, errors.New("invalid_access_token")
		}
		if id, err = strconv.Atoi(fmt.Sprintf("%v", destructID)); err != nil {
			return nil, errors.New("invalid_access_token")
		}
	}

	destructRoleID := c.RoleID
	if destructRoleID == "" {
		return nil, errors.New("invalid_access_token")
	}
	if rid, err = strconv.Atoi(fmt.Sprintf("%v", destructRoleID)); err != nil {
		if destructRoleID, err = aescrypt.DecryptAES(fmt.Sprintf("%v", destructRoleID), encryptionKey); err != nil {
			return nil, errors.New("invalid_access_token")
		}
		if rid, err = strconv.Atoi(fmt.Sprintf("%v", destructRoleID)); err != nil {
			return nil, errors.New("invalid_access_token")
		}
	}

	return &abstraction.AuthContext{
		ID:     id,
		RoleID: rid,
	}, nil
}

func (c AccessTokenClaims) RefreshTokenClaims() *RefreshTokenClaims {
	return &RefreshTokenClaims{
		ID:     c.ID,
		RoleID: c.RoleID,
		Exp:    time.Now().Add(config.JWT().RefreshTokenExpiry).Unix(),
	}
}

type RefreshTokenClaims struct {
	ID     string `json:"id"`
	RoleID string `json:"rid"`
	Exp    int64  `json:"exp"`

	jwt.RegisteredClaims
}

func (c RefreshTokenClaims) AuthContext() (*abstraction.AuthContext, error) {
	var (
		id  int
		rid int
		err error

		encryptionKey = config.JWT().EncryptionKey
	)

	destructID := c.ID
	if destructID == "" {
		return nil, errors.New("invalid_refresh_token")
	}
	if id, err = strconv.Atoi(fmt.Sprintf("%v", destructID)); err != nil {
		if destructID, err = aescrypt.DecryptAES(fmt.Sprintf("%v", destructID), encryptionKey); err != nil {
			return nil, errors.New("invalid_refresh_token")
		}
		if id, err = strconv.Atoi(fmt.Sprintf("%v", destructID)); err != nil {
			return nil, errors.New("invalid_refresh_token")
		}
	}

	destructRoleID := c.RoleID
	if destructRoleID == "" {
		return nil, errors.New("invalid_refresh_token")
	}
	if rid, err = strconv.Atoi(fmt.Sprintf("%v", destructRoleID)); err != nil {
		if destructRoleID, err = aescrypt.DecryptAES(fmt.Sprintf("%v", destructRoleID), encryptionKey); err != nil {
			return nil, errors.New("invalid_refresh_token")
		}
		if rid, err = strconv.Atoi(fmt.Sprintf("%v", destructRoleID)); err != nil {
			return nil, errors.New("invalid_refresh_token")
		}
	}

	return &abstraction.AuthContext{
		ID:     id,
		RoleID: rid,
	}, nil
}

func (c RefreshTokenClaims) AccessTokenClaims() *AccessTokenClaims {
	return &AccessTokenClaims{
		ID:     c.ID,
		RoleID: c.RoleID,
		Exp:    time.Now().Add(config.JWT().AccessTokenExpiry).Unix(),
	}
}
