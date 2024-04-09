package auth

import (
	"boilerplate/internal/abstraction"
	"boilerplate/internal/dto"
	"boilerplate/internal/factory"
	"boilerplate/pkg/util/response"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service Service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

// Login
// @Summary Login user
// @Description Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.AuthLoginRequest true "request body"
// @Success 200 {object} dto.AuthLoginResponseDoc
// @Failure 400 {object} response.ErrorResponse400
// @Failure 401 {object} response.ErrorResponse401
// @Failure 404 {object} response.ErrorResponse404
// @Failure 422 {object} response.ErrorResponse422
// @Failure 500 {object} response.ErrorResponse500
// @Router /auth/login [post]
func (h *handler) Login(c echo.Context) error {
	payload := new(dto.AuthLoginRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBadRequest(err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBadRequest(err).Send(c)
	}
	data, err := h.service.Login(c.(*abstraction.Context), payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(data).Send(c)
}

// Refresh Token
// @Summary Refresh Token user
// @Description Refresh Token user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "request body"
// @Success 200 {object} dto.RefreshTokenResponseDoc
// @Failure 400 {object} response.ErrorResponse400
// @Failure 401 {object} response.ErrorResponse401
// @Failure 404 {object} response.ErrorResponse404
// @Failure 422 {object} response.ErrorResponse422
// @Failure 500 {object} response.ErrorResponse500
// @Router /auth/refresh-token [post]
func (h *handler) RefreshToken(c echo.Context) error {
	payload := new(dto.RefreshTokenRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}
	data, err := h.service.RefreshToken(c.(*abstraction.Context), payload)
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(data).Send(c)
}

func (h *handler) Logout(c echo.Context) error {
	data, err := h.service.Logout(c.(*abstraction.Context))
	if err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(data).Send(c)
}
