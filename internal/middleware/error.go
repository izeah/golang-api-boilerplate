package middleware

import (
	"net/http"

	"boilerplate/pkg/util/response"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	var errCustom *response.Error

	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	switch report.Code {
	case http.StatusNotFound:
		errCustom = response.ErrorBuilder(&response.ErrorConstant.RouteNotFound, err)
	case http.StatusMethodNotAllowed:
		errCustom = response.ErrorBuilder(&response.ErrorConstant.MethodAllowedError, err)
	default:
		errCustom = response.ErrorBuilder(&response.ErrorConstant.InternalServerError, err)
	}

	if err = response.ErrorResponse(errCustom).Send(c); err != nil {
		c.Logger().Error(err)
	}
}
