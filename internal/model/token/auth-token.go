package modeltoken

import (
	"boilerplate/internal/config"

	"github.com/golang-jwt/jwt/v4"
)

type AuthToken struct {
	token *jwt.Token
}

func NewAuthToken(claims *AccessTokenClaims) *AuthToken {
	return &AuthToken{token: jwt.NewWithClaims(jwt.SigningMethodHS256, claims)}
}

func (t *AuthToken) AccessToken() (string, error) {
	signedString, err := t.token.SignedString([]byte(config.JWT().JWTKey))
	if err != nil {
		return "", err
	}
	return signedString, nil
}

func (t *AuthToken) RefreshToken() (string, error) {
	c := t.token.Claims.(*AccessTokenClaims)
	refreshTokenClaims := c.RefreshTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedString, err := token.SignedString([]byte(config.JWT().JWTRefKey))
	if err != nil {
		return "", err
	}
	return signedString, nil
}
