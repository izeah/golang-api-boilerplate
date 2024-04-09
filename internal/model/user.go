package model

import (
	"time"

	"boilerplate/internal/abstraction"
	"boilerplate/internal/config"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserEntity struct {
	Username     string `json:"username" validate:"required" example:"administrator"`
	Name         string `json:"name" validate:"required" example:"Lutfi Ramadhan"`
	Password     string `json:"password" validate:"required" gorm:"-" example:"nevemor3"`
	PasswordHash string `json:"-" gorm:"column:password"`
	Email        string `json:"email" validate:"required" example:"admin@console.code"`
	RoleID       int    `json:"role_id" required:"required" example:"1"`
	IsActive     *bool  `json:"is_active" validate:"required" gorm:"default:true" example:"true"`
}

// UserEntityModel ...
type UserEntityModel struct {
	// abstraction
	abstraction.Entity

	// entity
	UserEntity

	// context
	Context *abstraction.Context `json:"-" gorm:"-"`
}

// TableName ...
func (UserEntityModel) TableName() string {
	return "m_user"
}

func (m *UserEntityModel) BeforeCreate(_ *gorm.DB) (err error) {
	if m.Context != nil && m.Context.Auth != nil {
		m.CreatedBy = m.Context.Auth.ID
	}

	m.hashPassword()
	m.Password = ""
	return
}

func (m *UserEntityModel) BeforeUpdate(_ *gorm.DB) (err error) {
	if m.Password != "" {
		m.hashPassword()
		m.Password = ""
	}
	if m.Context != nil && m.Context.Auth != nil {
		m.ModifiedBy = &m.Context.Auth.ID
	}
	return
}

func (m *UserEntityModel) hashPassword() {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	m.PasswordHash = string(bytes)
}

func (m *UserEntityModel) GenerateToken() (string, error) {
	jwtKey := config.JWT().JWTKey

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  m.ID,
		"rid": m.RoleID,
		"exp": time.Now().Add(config.JWT().AccessTokenExpiry).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtKey))
	return tokenString, err
}
