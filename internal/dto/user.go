package dto

import (
	"strings"

	"boilerplate/internal/model"
	"boilerplate/pkg/util/response"

	"gorm.io/gorm"
)

// UserFilter ...
type UserFilter struct {
	ID       []int    `json:"id" query:"id"`
	Name     []string `json:"name" query:"name"`
	Username []string `json:"username" query:"username"`
	Email    []string `json:"email" query:"email"`
	RoleID   []int    `json:"role_id" query:"role_id"`
	IsActive *bool    `json:"is_active" query:"is_active"`

	Search *string `json:"search" query:"search"`
}

// Apply ...
func (f UserFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.Search != nil {
		db.Where("users.email ILIKE ? OR users.username ILIKE ? OR users.name ILIKE ?", "%"+*f.Search+"%", "%"+*f.Search+"%", "%"+*f.Search+"%")
	}

	if f.ID != nil {
		db.Where("id IN (?)", f.ID)
	}
	if f.Name != nil {
		names := []string{}
		for _, item := range f.Name {
			names = append(names, strings.ToLower(item))
		}
		db.Where("LOWER(name) SIMILAR TO ?", "%("+strings.Join(names, "|")+")%")
	}
	if f.Username != nil {
		usernames := []string{}
		for _, item := range f.Username {
			usernames = append(usernames, strings.ToLower(item))
		}
		db.Where("LOWER(username) SIMILAR TO ?", "%("+strings.Join(usernames, "|")+")%")
	}
	if f.Email != nil {
		emails := []string{}
		for _, item := range f.Email {
			emails = append(emails, strings.ToLower(item))
		}
		db.Where("LOWER(email) SIMILAR TO ?", "%("+strings.Join(emails, "|")+")%")
	}
	if f.RoleID != nil {
		db.Where("role_id IN (?)", f.RoleID)
	}
	if f.IsActive != nil {
		db.Where("is_active = ?", *f.IsActive)
	}
	return db
}

type FindUserResponseDoc struct {
	Meta response.Meta            `json:"meta"`
	Data []*model.UserEntityModel `json:"data"`
}

// UserFindByIDRequest ...
type UserFindByIDRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}

// UserFindByIDResponseDoc ...
type UserFindByIDResponseDoc struct {
	Meta response.Meta          `json:"meta"`
	Data *model.UserEntityModel `json:"data"`
}

// UserCreateRequest ...
type UserCreateRequest struct {
	Username string `form:"username" validate:"required" example:"administrator"`
	Name     string `form:"name" validate:"required" example:"Lutfi Ramadhan"`
	Password string `form:"password" validate:"required" gorm:"-" example:"nevemor3"`
	Email    string `form:"email" validate:"required" example:"admin@console.code"`
	RoleID   int    `form:"role_id" required:"required" example:"1"`
	IsActive *bool  `form:"is_active" validate:"required" example:"true"`
}

// UserCreateResponseDoc ...
type UserCreateResponseDoc struct {
	Meta response.Meta          `json:"meta"`
	Data *model.UserEntityModel `json:"data"`
}

// UserUpdateRequest ...
type UserUpdateRequest struct {
	ID       int    `param:"id" validate:"required,numeric"`
	Username string `form:"username" example:"administrator"`
	Name     string `form:"name" example:"Lutfi Ramadhan"`
	Password string `form:"password" gorm:"-" example:"nevemor3"`
	Email    string `form:"email" example:"admin@console.code"`
	RoleID   int    `form:"role_id" example:"1"`
	IsActive *bool  `form:"is_active" example:"true"`
}

// UserUpdateResponseDoc ...
type UserUpdateResponseDoc struct {
	Meta response.Meta          `json:"meta"`
	Data *model.UserEntityModel `json:"data"`
}

// UserDeleteRequest ...
type UserDeleteRequest struct {
	ID int `param:"id" validate:"required,numeric"`
}

// UserDeleteResponseDoc ...
type UserDeleteResponseDoc struct {
	Meta response.Meta `json:"meta"`
	Data interface{}   `json:"data"`
}
