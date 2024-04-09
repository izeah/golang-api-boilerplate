package dto

import (
	"fmt"

	"boilerplate/internal/config"
	"boilerplate/internal/model"
	modeltoken "boilerplate/internal/model/token"
	"boilerplate/pkg/util/response"

	"github.com/golang-jwt/jwt"
)

// AuthLoginRequest ...
type AuthLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// AuthLoginResponse ...
type AuthLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	model.UserEntityModel
}

// AuthLoginResponseDoc ...
type AuthLoginResponseDoc struct {
	Meta response.Meta     `json:"meta"`
	Data AuthLoginResponse `json:"data"`
}

// RefreshTokenRequest ...
type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AccessTokenClaims ...
func (r RefreshTokenRequest) AccessTokenClaims() (*modeltoken.AccessTokenClaims, error) {
	token, err := jwt.Parse(r.AccessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
		}
		return []byte(config.JWT().JWTKey), nil
	})
	if token == nil || !token.Valid || err != nil {
		if jwtErrValidation, ok := err.(*jwt.ValidationError); ok {
			c := token.Claims.(jwt.MapClaims)
			return &modeltoken.AccessTokenClaims{
				ID:     c["id"].(string),
				RoleID: c["rid"].(string),
				Exp:    int64(c["exp"].(float64)),
			}, jwtErrValidation
		}
		return nil, jwt.NewValidationError("invalid_access_token", jwt.ValidationErrorMalformed)
	}
	c := token.Claims.(jwt.MapClaims)
	return &modeltoken.AccessTokenClaims{
		ID:     c["id"].(string),
		RoleID: c["rid"].(string),
		Exp:    int64(c["exp"].(float64)),
	}, nil
}

// RefreshTokenClaims ...
func (r RefreshTokenRequest) RefreshTokenClaims() (*modeltoken.RefreshTokenClaims, error) {
	token, err := jwt.Parse(r.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
		}
		return []byte(config.JWT().JWTRefKey), nil
	})
	if token == nil || !token.Valid || err != nil {
		if jwtErrValidation, ok := err.(*jwt.ValidationError); ok {
			c := token.Claims.(jwt.MapClaims)
			return &modeltoken.RefreshTokenClaims{
				ID:     c["id"].(string),
				RoleID: c["rid"].(string),
				Exp:    int64(c["exp"].(float64)),
			}, jwtErrValidation
		}
		return nil, jwt.NewValidationError("invalid_refresh_token", jwt.ValidationErrorMalformed)
	}
	c := token.Claims.(jwt.MapClaims)
	return &modeltoken.RefreshTokenClaims{
		ID:     c["id"].(string),
		RoleID: c["rid"].(string),
		Exp:    int64(c["exp"].(float64)),
	}, nil
}

// RefreshTokenResponse ...
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponseDoc ...
type RefreshTokenResponseDoc struct {
	Meta response.Meta        `json:"meta"`
	Data RefreshTokenResponse `json:"data"`
}
