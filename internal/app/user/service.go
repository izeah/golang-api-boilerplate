package user

import (
	"errors"
	"fmt"
	"math"
	"net/http"

	"boilerplate/internal/abstraction"
	"boilerplate/internal/dto"
	"boilerplate/internal/factory"
	"boilerplate/internal/model"
	"boilerplate/internal/repository"
	"boilerplate/pkg/util/response"
	"boilerplate/pkg/util/trxmanager"

	"gorm.io/gorm"
)

type Service interface {
	Find(ctx *abstraction.Context, f *dto.UserFilter, p *abstraction.Pagination) ([]*model.UserEntityModel, *abstraction.PaginationInfo, error)
	FindByID(ctx *abstraction.Context, payload *dto.UserFindByIDRequest) (*model.UserEntityModel, error)
	Create(ctx *abstraction.Context, payload *dto.UserCreateRequest) (*model.UserEntityModel, error)
	Update(ctx *abstraction.Context, payload *dto.UserUpdateRequest) (*model.UserEntityModel, error)
	Delete(ctx *abstraction.Context, payload *dto.UserDeleteRequest) error
}

type service struct {
	UserRepository repository.User

	DB *gorm.DB
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepository: f.UserRepository,

		DB: f.DB,
	}
}

func (s *service) Find(ctx *abstraction.Context, f *dto.UserFilter, p *abstraction.Pagination) (data []*model.UserEntityModel, info *abstraction.PaginationInfo, err error) {
	if data, info, err = s.UserRepository.Find(ctx, f, p); err != nil {
		return nil, nil, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}
	if p != nil && p.PageSize != nil {
		info.Pages = int(math.Ceil(float64(info.Count) / float64(*p.PageSize)))
		if len(data) > *p.PageSize {
			data = data[:len(data)-1]
			info.MoreRecords = true
		}
	}
	return
}

func (s *service) FindByID(ctx *abstraction.Context, payload *dto.UserFindByIDRequest) (data *model.UserEntityModel, err error) {
	if data, err = s.UserRepository.FindByID(ctx, payload.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
		}
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}
	return
}

func (s *service) Create(ctx *abstraction.Context, payload *dto.UserCreateRequest) (data *model.UserEntityModel, err error) {
	if data, err = s.UserRepository.FindByUsernameOrEmail(ctx, payload.Username, payload.Email); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}
	if data.ID != 0 {
		return nil, response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, fmt.Sprintf("Email %s or username %s already exist", payload.Email, payload.Username))
	}
	if err = trxmanager.New(s.DB).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data = &model.UserEntityModel{}
		data.Context = ctx
		data.UserEntity = model.UserEntity{
			Name:     payload.Name,
			Email:    payload.Email,
			RoleID:   payload.RoleID,
			Username: payload.Username,
			Password: payload.Password,
			IsActive: payload.IsActive,
		}
		if err = s.UserRepository.Create(ctx, &data).Error; err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return
}

func (s *service) Update(ctx *abstraction.Context, payload *dto.UserUpdateRequest) (data *model.UserEntityModel, err error) {
	if data, err = s.UserRepository.FindByID(ctx, payload.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
		}
		return nil, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}
	if data.Username != payload.Username || data.Email != payload.Email {
		if data, err = s.UserRepository.FindByUsernameOrEmail(ctx, payload.Username, payload.Email); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		if data.ID != 0 {
			return nil, response.CustomErrorBuilder(http.StatusBadRequest, response.E_BAD_REQUEST, fmt.Sprintf("Email %s or username %s already exist", payload.Email, payload.Username))
		}
	}
	if err = trxmanager.New(s.DB).WithTrx(ctx, func(ctx *abstraction.Context) error {
		data = &model.UserEntityModel{}
		data.Context = ctx
		data.ID = payload.ID
		data.UserEntity = model.UserEntity{
			Name:     payload.Name,
			Email:    payload.Email,
			RoleID:   payload.RoleID,
			Username: payload.Username,
			Password: payload.Password,
			IsActive: payload.IsActive,
		}
		if err = s.UserRepository.Update(ctx, data).Error; err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return
}

func (s *service) Delete(ctx *abstraction.Context, payload *dto.UserDeleteRequest) error {
	if _, err := s.UserRepository.FindByID(ctx, payload.ID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ErrorBuilder(&response.ErrorConstant.NotFound, err)
		}
		return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
	}
	return trxmanager.New(s.DB).WithTrx(ctx, func(ctx *abstraction.Context) error {
		if err := s.UserRepository.Delete(ctx, &dto.UserFilter{ID: []int{payload.ID}}).Error; err != nil {
			return response.ErrorBuilder(&response.ErrorConstant.UnprocessableEntity, err)
		}
		return nil
	})
}
