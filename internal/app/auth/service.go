package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"boilerplate/internal/abstraction"
	"boilerplate/internal/config"
	"boilerplate/internal/dto"
	"boilerplate/internal/factory"
	modeltoken "boilerplate/internal/model/token"
	"boilerplate/internal/repository"
	"boilerplate/pkg/redis"
	"boilerplate/pkg/util/aescrypt"
	"boilerplate/pkg/util/response"

	"github.com/golang-jwt/jwt"
	goRedis "github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service interface {
	Login(ctx *abstraction.Context, payload *dto.AuthLoginRequest) (*dto.AuthLoginResponse, error)
	RefreshToken(ctx *abstraction.Context, payload *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	Logout(ctx *abstraction.Context) (map[string]interface{}, error)
}

type service struct {
	UserRepository repository.User

	DB *gorm.DB
}

func NewService(f *factory.Factory) *service {
	return &service{
		UserRepository: f.UserRepository,

		DB: f.DB,
	}
}

func (s *service) Login(ctx *abstraction.Context, payload *dto.AuthLoginRequest) (*dto.AuthLoginResponse, error) {
	data, err := s.UserRepository.FindByUsernameOrEmail(ctx, payload.Username, payload.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
		}
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}

	if data.IsActive == nil || (data.IsActive != nil && !*data.IsActive) {
		return nil, response.ErrorBuilder(&response.ErrorConstant.Unauthorized, errors.New("account is not active"))
	}

	if err = bcrypt.CompareHashAndPassword([]byte(data.PasswordHash), []byte(payload.Password)); err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.Unauthorized, errors.New("password is incorrect"))
	}

	authLoggedInUser, err := redis.Client().Get(ctx.Request().Context(), fmt.Sprintf("auth_user_id_%d_flag", data.ID)).Result()
	if err != nil && !errors.Is(err, goRedis.Nil) {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}
	if authLoggedInUser == "1" {
		authLoggedInUserInfo, err := redis.Client().HGetAll(ctx.Request().Context(), fmt.Sprintf("auth_user_id_%d_info", data.ID)).Result()
		if err != nil && !errors.Is(err, goRedis.Nil) {
			return nil, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		if authLoggedInUserInfo["ip_address"] != ctx.RealIP() || authLoggedInUserInfo["user_agent"] != ctx.Request().UserAgent() {
			return nil, response.ErrorBuilder(&response.ErrorConstant.UnauthorizedNewDevice, errors.New("this_user_already_logged_in"))
		}
	}

	var encryptedUserID, encryptedRoleID string
	if encryptedUserID, err = s.encryptTokenClaims(data.ID); err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}
	if encryptedRoleID, err = s.encryptTokenClaims(data.RoleID); err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}

	accessTokenClaims := &modeltoken.AccessTokenClaims{
		ID:     encryptedUserID,
		RoleID: encryptedRoleID,
		Exp:    time.Now().Add(config.JWT().AccessTokenExpiry).Unix(),
	}
	authToken := modeltoken.NewAuthToken(accessTokenClaims)
	accessToken, err := authToken.AccessToken()
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}
	refreshToken, err := authToken.RefreshToken()
	if err != nil {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}

	_ = redis.Client().Set(ctx.Request().Context(), fmt.Sprintf("auth_user_id_%d_flag", data.ID), true, config.JWT().AccessTokenExpiry)
	_ = redis.Client().HSet(ctx.Request().Context(), fmt.Sprintf("auth_user_id_%d_info", data.ID), map[string]interface{}{
		"ip_address": ctx.RealIP(),
		"user_agent": ctx.Request().UserAgent(),
		"token":      accessToken,
	})
	_ = redis.Client().Expire(ctx.Request().Context(), fmt.Sprintf("auth_user_id_%d_info", data.ID), config.JWT().AccessTokenExpiry)

	return &dto.AuthLoginResponse{
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		UserEntityModel: *data,
	}, nil
}

func (s *service) RefreshToken(ctx *abstraction.Context, payload *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	accessTokenClaims, err := payload.AccessTokenClaims()
	if err != nil && err.(*jwt.ValidationError).Errors != jwt.ValidationErrorExpired {
		return nil, response.CustomErrorBuilder(http.StatusBadRequest, "invalid_access_token", "invalid_access_token")
	}
	accessTokenAuthCtx, err := accessTokenClaims.AuthContext()
	if err != nil {
		return nil, response.CustomErrorBuilder(http.StatusBadRequest, err.Error(), "invalid_access_token")
	}

	refreshTokenClaims, err := payload.RefreshTokenClaims()
	if err != nil {
		if jwtValErr := err.(*jwt.ValidationError); jwtValErr.Errors == jwt.ValidationErrorExpired {
			return nil, response.CustomErrorBuilder(http.StatusUnauthorized, "refresh_token_is_expired", "refresh_token_is_expired")
		} else {
			return nil, response.CustomErrorBuilder(http.StatusBadRequest, jwtValErr.Error(), "invalid_refresh_token")
		}
	}
	refreshTokenAuthCtx, err := refreshTokenClaims.AuthContext()
	if err != nil {
		return nil, response.CustomErrorBuilder(http.StatusBadRequest, err.Error(), "invalid_refresh_token")
	}

	if refreshTokenAuthCtx.ID != accessTokenAuthCtx.ID || refreshTokenAuthCtx.RoleID != accessTokenAuthCtx.RoleID {
		return nil, response.CustomErrorBuilder(http.StatusUnauthorized, "unauthorized_to_refresh_token", "unauthorized_to_refresh_token")
	}

	accessTokenClaims = refreshTokenClaims.AccessTokenClaims()
	authToken := modeltoken.NewAuthToken(accessTokenClaims)
	accessToken, err := authToken.AccessToken()
	if err != nil {
		return nil, response.CustomErrorBuilder(http.StatusUnauthorized, err.Error(), "err_generate_access_token")
	}
	refreshToken, err := authToken.RefreshToken()
	if err != nil {
		return nil, response.CustomErrorBuilder(http.StatusUnauthorized, err.Error(), "err_generate_refresh_token")
	}

	_ = redis.Client().Del(ctx.Request().Context(), fmt.Sprintf("auth_user_id_%d_flag", accessTokenAuthCtx.ID))
	_ = redis.Client().Del(ctx.Request().Context(), fmt.Sprintf("auth_user_id_%d_info", accessTokenAuthCtx.ID))
	_ = redis.Client().Set(ctx.Request().Context(), fmt.Sprintf("auth_user_id_%d_flag", accessTokenAuthCtx.ID), true, config.JWT().AccessTokenExpiry)
	_ = redis.Client().HSet(ctx.Request().Context(), fmt.Sprintf("auth_user_id_%d_info", accessTokenAuthCtx.ID), map[string]interface{}{
		"ip_address": ctx.RealIP(),
		"user_agent": ctx.Request().UserAgent(),
		"token":      accessToken,
	})
	_ = redis.Client().Expire(ctx.Request().Context(), fmt.Sprintf("auth_user_id_%d_info", accessTokenAuthCtx.ID), config.JWT().AccessTokenExpiry)

	return &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *service) encryptTokenClaims(v int) (encryptedString string, err error) {
	encryptedString, err = aescrypt.EncryptAES(fmt.Sprint(v), config.JWT().EncryptionKey)
	return
}

func (s *service) Logout(ctx *abstraction.Context) (map[string]interface{}, error) {
	tokenString := strings.Replace(ctx.Request().Header.Get("Authorization"), "Bearer ", "", -1)

	result, err := redis.Client().HGetAll(ctx.Request().Context(), fmt.Sprintf("auth_user_id_%d_info", ctx.Auth.ID)).Result()
	if err != nil && !errors.Is(err, goRedis.Nil) {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}

	if result == nil || result["token"] == "" || result["token"] != tokenString {
		return map[string]interface{}{
			"message": "Another user has logged out this account",
		}, nil
	}

	_ = redis.Client().Del(ctx.Request().Context(), fmt.Sprintf("auth_user_id_%d_flag", ctx.Auth.ID))
	_ = redis.Client().Del(ctx.Request().Context(), fmt.Sprintf("auth_user_id_%d_info", ctx.Auth.ID))

	return map[string]interface{}{
		"message": "Logout successful",
	}, nil
}
