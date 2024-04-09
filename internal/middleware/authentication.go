package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"boilerplate/internal/abstraction"
	"boilerplate/internal/config"
	"boilerplate/pkg/redis"
	"boilerplate/pkg/util/aescrypt"
	"boilerplate/pkg/util/response"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id, rid          int
			jwtKey           = config.JWT().JWTKey
			jwtEncryptionKey = config.JWT().EncryptionKey
		)

		authToken := c.Request().Header.Get("Authorization")
		if authToken == "" {
			return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
		}
		if !strings.Contains(authToken, "Bearer") {
			return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
		}

		tokenString := strings.Replace(authToken, "Bearer ", "", -1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
			}
			return []byte(jwtKey), nil
		})
		if token == nil || !token.Valid || err != nil {
			if errJWT, ok := err.(*jwt.ValidationError); ok {
				if errJWT.Errors == jwt.ValidationErrorExpired {
					destructID := token.Claims.(jwt.MapClaims)["id"]
					if destructID == nil {
						return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
					}
					if id, err = strconv.Atoi(fmt.Sprintf("%v", destructID)); err != nil {
						if destructID, err = aescrypt.DecryptAES(fmt.Sprintf("%v", destructID), jwtEncryptionKey); err != nil {
							return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
						}
						if id, err = strconv.Atoi(fmt.Sprintf("%v", destructID)); err != nil {
							return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
						}
					}
					_ = redis.Client().Del(c.Request().Context(), fmt.Sprintf("auth_user_id_%d_flag", id))
					_ = redis.Client().Del(c.Request().Context(), fmt.Sprintf("auth_user_id_%d_info", id))
					return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "access_token_is_expired").Send(c)
				}
				return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, errJWT.Error()).Send(c)
			}
			return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return response.ErrorBuilder(&response.ErrorConstant.Unauthorized, err).Send(c)
		}

		destructID := claims["id"]
		if destructID == nil {
			return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
		}
		if id, err = strconv.Atoi(fmt.Sprintf("%v", destructID)); err != nil {
			if destructID, err = aescrypt.DecryptAES(fmt.Sprintf("%v", destructID), jwtEncryptionKey); err != nil {
				return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
			}
			if id, err = strconv.Atoi(fmt.Sprintf("%v", destructID)); err != nil {
				return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
			}
		}

		destructRoleID := claims["rid"]
		if destructRoleID == nil {
			return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
		}
		if rid, err = strconv.Atoi(fmt.Sprintf("%v", destructRoleID)); err != nil {
			if destructRoleID, err = aescrypt.DecryptAES(fmt.Sprintf("%v", destructRoleID), jwtEncryptionKey); err != nil {
				return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
			}
			if rid, err = strconv.Atoi(fmt.Sprintf("%v", destructRoleID)); err != nil {
				return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
			}
		}

		cc := c.(*abstraction.Context)
		cc.Auth = &abstraction.AuthContext{
			ID:     id,
			RoleID: rid,
		}

		return next(cc)
	}
}

func Logout(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			id, rid          int
			jwtKey           = config.JWT().JWTKey
			jwtEncryptionKey = config.JWT().EncryptionKey
		)

		authToken := c.Request().Header.Get("Authorization")
		if authToken == "" {
			return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
		}
		if !strings.Contains(authToken, "Bearer") {
			return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
		}

		tokenString := strings.Replace(authToken, "Bearer ", "", -1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method :%v", token.Header["alg"])
			}
			return []byte(jwtKey), nil
		})
		if token == nil || !token.Valid || err != nil {
			if errJWT, ok := err.(*jwt.ValidationError); ok {
				if errJWT.Errors != jwt.ValidationErrorExpired {
					return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, errJWT.Error()).Send(c)
				}
			} else {
				return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
			}
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return response.ErrorBuilder(&response.ErrorConstant.Unauthorized, err).Send(c)
		}

		destructID := claims["id"]
		if destructID == nil {
			return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
		}
		if id, err = strconv.Atoi(fmt.Sprintf("%v", destructID)); err != nil {
			if destructID, err = aescrypt.DecryptAES(fmt.Sprintf("%v", destructID), jwtEncryptionKey); err != nil {
				return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
			}
			if id, err = strconv.Atoi(fmt.Sprintf("%v", destructID)); err != nil {
				return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
			}
		}

		destructRoleID := claims["rid"]
		if destructRoleID == nil {
			return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
		}
		if rid, err = strconv.Atoi(fmt.Sprintf("%v", destructRoleID)); err != nil {
			if destructRoleID, err = aescrypt.DecryptAES(fmt.Sprintf("%v", destructRoleID), jwtEncryptionKey); err != nil {
				return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
			}
			if rid, err = strconv.Atoi(fmt.Sprintf("%v", destructRoleID)); err != nil {
				return response.CustomErrorBuilder(http.StatusUnauthorized, response.E_UNAUTHORIZED, "invalid_token").Send(c)
			}
		}

		cc := c.(*abstraction.Context)
		cc.Auth = &abstraction.AuthContext{
			ID:     id,
			RoleID: rid,
		}

		return next(cc)
	}
}
