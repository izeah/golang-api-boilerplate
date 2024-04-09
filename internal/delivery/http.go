package delivery

import (
	"fmt"
	"net/http"

	"boilerplate/internal/app/user"
	"boilerplate/internal/config"
	"boilerplate/internal/factory"
	"boilerplate/pkg/database"

	"github.com/labstack/echo/v4"

	_docs "boilerplate/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func HTTP(e *echo.Echo, f *factory.Factory) {
	var (
		APP     = config.App().Name
		VERSION = config.App().Version
		HOST    = config.App().Host
		SCHEME  = config.App().Schemes
		ENV     = config.App().ENV
	)

	// index
	e.GET("/", func(c echo.Context) error {
		message := fmt.Sprintf("Welcome to %s version %s %s", APP, VERSION, ENV)
		return c.String(http.StatusOK, message)
	})

	// doc
	_docs.SwaggerInfo.Title = APP
	_docs.SwaggerInfo.Version = VERSION
	_docs.SwaggerInfo.Host = HOST
	_docs.SwaggerInfo.BasePath = ""
	_docs.SwaggerInfo.Schemes = append(SCHEME, "https")
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	user.NewHandler(f).Route(e.Group("/user"))

	e.GET("/position", func(c echo.Context) error {
		var data map[string]any
		if err := database.MYSQL().WithContext(c.Request().Context()).Table("tc_positions").Limit(1).Find(&data).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, nil)
		}
		return c.JSON(http.StatusOK, data)
	})
}
