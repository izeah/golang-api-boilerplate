package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"boilerplate/internal/config"
	"boilerplate/pkg/util/validator"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func Init(e *echo.Echo) {
	NAME := fmt.Sprintf("%s-%s", config.App().Name, config.App().ENV)

	e.Use(Context)
	e.Use(
		echoMiddleware.Recover(),
		echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "apikey"},
			AllowMethods: []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodDelete},
		}),
		echoMiddleware.LoggerWithConfig(echoMiddleware.LoggerConfig{
			Format:           fmt.Sprintf("\n%s | ${host} | ${time_custom} | ${status} | ${latency_human} | ${remote_ip} | ${method} | ${uri}", NAME),
			CustomTimeFormat: "2006/01/02 15:04:05",
			Output:           os.Stdout,
		}),
		echoMiddleware.TimeoutWithConfig(echoMiddleware.TimeoutConfig{
			Skipper:      echoMiddleware.DefaultSkipper,
			ErrorMessage: http.StatusText(http.StatusRequestTimeout),
			Timeout:      5 * time.Minute,
		}),
	)
	e.HTTPErrorHandler = ErrorHandler
	e.Validator = &validator.CustomValidator{Validator: validator.NewValidator()}
}
