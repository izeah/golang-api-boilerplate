package user

import (
	"boilerplate/internal/middleware"

	"github.com/labstack/echo/v4"
)

// Route ...
func (h *handler) Route(v *echo.Group) {
	v.GET("", h.Find, middleware.Authentication)
	v.GET("/:id", h.FindByID, middleware.Authentication)
	v.POST("", h.Create, middleware.Authentication)
	v.PUT("/:id", h.Update, middleware.Authentication)
	v.DELETE("/:id", h.Delete, middleware.Authentication)
}
