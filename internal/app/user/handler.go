package user

import (
	"boilerplate/internal/abstraction"
	"boilerplate/internal/dto"
	"boilerplate/internal/factory"
	"boilerplate/internal/model"
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

// Find User
// @Summary Find User
// @Description Find User
// @Tags User
// @Produce json
// @Param request query dto.UserFilter true "request query"
// @param request query abstraction.Pagination true "request query pagination"
// @Success 200 {object} dto.FindUserResponseDoc
// @Failure 400 {object} response.ErrorResponse400
// @Failure 401 {object} response.ErrorResponse401
// @Failure 404 {object} response.ErrorResponse404
// @Failure 422 {object} response.ErrorResponse422
// @Failure 500 {object} response.ErrorResponse500
// @Router /user [get]
func (h *handler) Find(c echo.Context) (err error) {
	f := new(dto.UserFilter)
	if err := c.Bind(f); err != nil {
		return response.ErrorBadRequest(err).Send(c)
	}

	p := new(abstraction.Pagination)
	if err := c.Bind(p); err != nil {
		return response.ErrorBadRequest(err).Send(c)
	}

	var (
		data []*model.UserEntityModel
		info *abstraction.PaginationInfo
	)
	if data, info, err = h.service.Find(c.(*abstraction.Context), f, p); err != nil {
		return response.ErrorResponse(err).Send(c)
	}

	return response.SuccessResponse(data).WithPagination(info).Send(c)
}

// Find User By ID
// @Summary Find User by ID
// @Description Find User by ID
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {object} dto.UserFindByIDResponseDoc
// @Failure 400 {object} response.ErrorResponse400
// @Failure 401 {object} response.ErrorResponse401
// @Failure 404 {object} response.ErrorResponse404
// @Failure 422 {object} response.ErrorResponse422
// @Failure 500 {object} response.ErrorResponse500
// @Router /user/{id} [get]
func (h *handler) FindByID(c echo.Context) (err error) {
	payload := new(dto.UserFindByIDRequest)
	if err = c.Bind(payload); err != nil {
		return response.ErrorBadRequest(err).Send(c)
	}
	if err = c.Validate(payload); err != nil {
		return response.ErrorBadRequest(err).Send(c)
	}
	var data *model.UserEntityModel
	if data, err = h.service.FindByID(c.(*abstraction.Context), payload); err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(data).Send(c)
}

// Create User
// @Summary Create User
// @Description Create User
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UserCreateRequest true "request body"
// @Success 200 {object} dto.UserCreateResponseDoc
// @Failure 400 {object} response.ErrorResponse400
// @Failure 401 {object} response.ErrorResponse401
// @Failure 404 {object} response.ErrorResponse404
// @Failure 422 {object} response.ErrorResponse422
// @Failure 500 {object} response.ErrorResponse500
// @Router /user [post]
func (h *handler) Create(c echo.Context) (err error) {
	payload := new(dto.UserCreateRequest)
	if err = c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err = c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}
	var data *model.UserEntityModel
	if data, err = h.service.Create(c.(*abstraction.Context), payload); err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(data).Send(c)
}

// Update godoc
// @Summary Update User
// @Description Update User
// @Tags User
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param id path int true "id path"
// @Param request body dto.UserUpdateRequest true "request body"
// @Success 200 {object} dto.UserUpdateResponseDoc
// @Failure 400 {object} response.ErrorResponse400
// @Failure 401 {object} response.ErrorResponse401
// @Failure 404 {object} response.ErrorResponse404
// @Failure 422 {object} response.ErrorResponse422
// @Failure 500 {object} response.ErrorResponse500
// @Router /user/{id} [put]
func (h *handler) Update(c echo.Context) (err error) {
	payload := new(dto.UserUpdateRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.BadRequest, err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBuilder(&response.ErrorConstant.Validation, err).Send(c)
	}
	var data *model.UserEntityModel
	if data, err = h.service.Update(c.(*abstraction.Context), payload); err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(data).Send(c)
}

// Delete User
// @Summary Delete User
// @Description Delete User
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "id path"
// @Success 200 {object} dto.UserDeleteResponseDoc
// @Failure 400 {object} response.ErrorResponse400
// @Failure 401 {object} response.ErrorResponse401
// @Failure 404 {object} response.ErrorResponse404
// @Failure 422 {object} response.ErrorResponse422
// @Failure 500 {object} response.ErrorResponse500
// @Router /user/{id} [delete]
func (h *handler) Delete(c echo.Context) error {
	payload := new(dto.UserDeleteRequest)
	if err := c.Bind(payload); err != nil {
		return response.ErrorBadRequest(err).Send(c)
	}
	if err := c.Validate(payload); err != nil {
		return response.ErrorBadRequest(err).Send(c)
	}
	if err := h.service.Delete(c.(*abstraction.Context), payload); err != nil {
		return response.ErrorResponse(err).Send(c)
	}
	return response.SuccessResponse(nil).Send(c)
}
