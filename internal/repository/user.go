package repository

import (
	"errors"

	"boilerplate/internal/abstraction"
	"boilerplate/internal/dto"
	"boilerplate/internal/model"

	"gorm.io/gorm"
)

type User interface {
	Find(ctx *abstraction.Context, f *dto.UserFilter, p *abstraction.Pagination) ([]*model.UserEntityModel, *abstraction.PaginationInfo, error)
	FindByUsernameOrEmail(ctx *abstraction.Context, username, email string) (data *model.UserEntityModel, err error)
	FindByID(ctx *abstraction.Context, id int) (data *model.UserEntityModel, err error)
	Create(ctx *abstraction.Context, e interface{}) *gorm.DB
	Update(ctx *abstraction.Context, e *model.UserEntityModel) *gorm.DB
	Delete(ctx *abstraction.Context, f *dto.UserFilter) *gorm.DB
}

type user struct {
	abstraction.Repository
}

func NewUser(db *gorm.DB) User {
	return &user{
		Repository: abstraction.Repository{
			Db: db,
		},
	}
}

func (r *user) Find(ctx *abstraction.Context, f *dto.UserFilter, p *abstraction.Pagination) ([]*model.UserEntityModel, *abstraction.PaginationInfo, error) {
	var (
		data  []*model.UserEntityModel
		count int64
		err   error

		info = &abstraction.PaginationInfo{Pagination: p}
	)

	if err = r.CheckTrx(ctx).Model(&model.UserEntityModel{}).Scopes(func(db *gorm.DB) *gorm.DB {
		if f != nil {
			f.Apply(db)
		}
		return db
	}).Count(&count).Error; err != nil {
		return nil, nil, err
	}

	if err = r.CheckTrx(ctx).Model(&model.UserEntityModel{}).Scopes(func(db *gorm.DB) *gorm.DB {
		if f != nil {
			f.Apply(db)
		}
		if p != nil {
			if p.Page == nil || p.PageSize == nil {
				p.Init()
			}
			return db.Offset(p.GetOffset()).Limit(p.GetLimit()).Order(p.GetOrderBy())
		}
		return db
	}).Find(&data).Error; err != nil {
		return nil, nil, err
	}

	info.Count = count
	return data, info, nil
}

func (r *user) FindByUsernameOrEmail(ctx *abstraction.Context, username, email string) (data *model.UserEntityModel, err error) {
	err = r.CheckTrx(ctx).Where("username = ? OR email = ?", username, email).Take(&data).Error
	return
}

func (r *user) FindByID(ctx *abstraction.Context, id int) (data *model.UserEntityModel, err error) {
	err = r.CheckTrx(ctx).Where("id = ?", id).Take(&data).Error
	return
}

func (r *user) Create(ctx *abstraction.Context, e interface{}) *gorm.DB {
	return r.CheckTrx(ctx).Create(e)
}

func (r *user) Update(ctx *abstraction.Context, e *model.UserEntityModel) *gorm.DB {
	return r.CheckTrx(ctx).Save(e)
}

func (r *user) Delete(ctx *abstraction.Context, f *dto.UserFilter) *gorm.DB {
	return r.CheckTrx(ctx).Scopes(func(db *gorm.DB) *gorm.DB {
		if f == nil {
			db.AddError(errors.New("need filter"))
		}
		return f.Apply(db)
	}).Delete(&model.UserEntityModel{})
}
