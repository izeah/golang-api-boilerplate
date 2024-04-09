package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Meta        Meta        `json:"meta"`
	Error       interface{} `json:"data"`
	Description interface{} `json:"description,omitempty"`
}

type Error struct {
	Header       *http.Header
	Response     errorResponse `json:"response"`
	Code         int           `json:"code"`
	ErrorMessage error
}

const (
	E_DUPLICATE            = "duplicate"
	E_NOT_FOUND            = "not_found"
	E_UNPROCESSABLE_ENTITY = "unprocessable_entity"
	E_UNAUTHORIZED         = "unauthorized"
	E_FORBIDDEN            = "forbidden"
	E_METHOD_NOT_ALLOWED   = "method_not_allowed"
	E_BAD_REQUEST          = "bad_request"
	E_SERVER_ERROR         = "server_error"
	E_TOO_MANY_REQUEST     = "too_many_request"
)

type errorConstant struct {
	Duplicate               Error
	NotFound                Error
	RouteNotFound           Error
	UnprocessableEntity     Error
	Unauthorized            Error
	UnauthorizedNewDevice   Error
	BadRequest              Error
	Forbidden               Error
	Validation              Error
	MethodAllowedError      Error
	InternalServerError     Error
	ServiceUnavailableError Error
	NotFileUpload           Error
	UploadFileError         Error

	TooManyRequest func(retryAfterSecond float64) *Error
}

var (
	_validator    = validator.New()
	ErrorConstant = errorConstant{
		ServiceUnavailableError: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Service Unavailable",
				},
				Error: E_SERVER_ERROR,
			},
			Code: http.StatusServiceUnavailable,
		},
		MethodAllowedError: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Method Not Allowed",
				},
				Error: E_METHOD_NOT_ALLOWED,
			},
			Code: http.StatusMethodNotAllowed,
		},
		Duplicate: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Created value already exists",
				},
				Error: E_DUPLICATE,
			},
			Code: http.StatusConflict,
		},
		NotFound: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Data not found",
				},
				Error: E_NOT_FOUND,
			},
			Code: http.StatusNotFound,
		},
		RouteNotFound: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Route not found",
				},
				Error: E_NOT_FOUND,
			},
			Code: http.StatusNotFound,
		},
		UnprocessableEntity: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Invalid parameters or payload",
				},
				Error: E_UNPROCESSABLE_ENTITY,
			},
			Code: http.StatusUnprocessableEntity,
		},
		Unauthorized: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Unauthorized, please login",
				},
				Error: E_UNAUTHORIZED,
			},
			Code: http.StatusUnauthorized,
		},
		UnauthorizedNewDevice: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Unauthorized, this account already login with old device",
				},
				Error: E_UNAUTHORIZED,
			},
			Code: http.StatusUnauthorized,
		},
		Forbidden: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Forbidden access",
				},
				Error: E_FORBIDDEN,
			},
			Code: http.StatusForbidden,
		},
		BadRequest: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Bad Request",
				},
				Error: E_BAD_REQUEST,
			},
			Code: http.StatusBadRequest,
		},
		Validation: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Invalid parameters or payload",
				},
				Error: E_BAD_REQUEST,
			},
			Code: http.StatusBadRequest,
		},
		InternalServerError: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Something bad happened",
				},
				Error: E_SERVER_ERROR,
			},
			Code: http.StatusInternalServerError,
		},
		NotFileUpload: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "No files to upload",
				},
				Error: E_SERVER_ERROR,
			},
			Code: http.StatusInternalServerError,
		},
		TooManyRequest: func(retryAfterSecond float64) *Error {
			if retryAfterSecond = math.Round(retryAfterSecond); retryAfterSecond == 0 {
				retryAfterSecond = 1
			}
			return &Error{
				Header: &http.Header{echo.HeaderRetryAfter: []string{strconv.Itoa(int(retryAfterSecond))}},
				Response: errorResponse{
					Meta: Meta{
						Success: false,
						Message: "Too Many Request",
					},
					Error: E_TOO_MANY_REQUEST,
				},
				Code: http.StatusTooManyRequests,
			}
		},
		UploadFileError: Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Failed to upload file",
				},
				Error: E_SERVER_ERROR,
			},
			Code: http.StatusInternalServerError,
		},
	}
)

func ErrorBuilder(res *Error, message error, vals ...interface{}) *Error {
	res.ErrorMessage = message
	res.Response.Description = vals
	return res
}

// ErrorBadRequest ...
func ErrorBadRequest(err error) *Error {
	if err == nil {
		return nil
	}

	var eR *Error
	if errors.As(err, &eR) {
		return eR
	}

	var e *echo.HTTPError
	if errors.As(err, &e) {
		var (
			data                map[string]interface{}
			message             = "bad request"
			eUnmarshalTypeError *json.UnmarshalTypeError
		)
		if errors.As(e.Internal, &eUnmarshalTypeError) {
			data = make(map[string]interface{})
			data[eUnmarshalTypeError.Field] = map[string]interface{}{
				"message": fmt.Sprintf("should be %s!", eUnmarshalTypeError.Type),
			}
		} else {
			var eSyntaxError *json.SyntaxError
			if errors.As(e.Internal, &eSyntaxError) {
				message = eSyntaxError.Error()
			}
		}
		return &Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: message,
				},
				Description: data,
			},
			Code: http.StatusBadRequest,
		}
	}

	if err = _validator.Struct(err); e != nil {
		var errValidation validator.ValidationErrors
		if errors.As(err, &errValidation) {
			var (
				data = make(map[string]interface{})
				ef   validator.FieldError
			)
			for _, ef = range errValidation {
				data[strings.ToLower(ef.Field())] = map[string]string{
					"tag":     ef.Tag(),
					"param":   ef.Param(),
					"message": Translate(ef),
				}
			}
			return &Error{
				Response: errorResponse{
					Meta: Meta{
						Success: false,
						Message: "Bad Request",
					},
					Description: data,
				},
				Code: http.StatusBadRequest,
			}
		}
	}

	return ErrorBuilder(&ErrorConstant.BadRequest, err, err.Error())
}

func CustomErrorBuilder(code int, err interface{}, message string, vals ...interface{}) *Error {
	return &Error{
		Response: errorResponse{
			Meta: Meta{
				Success: false,
				Message: message,
			},
			Error:       err,
			Description: vals,
		},
		Code:         code,
		ErrorMessage: errors.New(message),
	}
}

// RestyErrorBuilder ...
func RestyErrorBuilder(resp *resty.Response, meta Meta) error {
	if resp.StatusCode() >= http.StatusMultipleChoices || resp.StatusCode() < http.StatusOK || !meta.Success {
		switch resp.StatusCode() {
		case http.StatusNotFound:
			return &ErrorConstant.NotFound
		case http.StatusConflict:
			return &ErrorConstant.Duplicate
		case http.StatusUnauthorized:
			return &ErrorConstant.Unauthorized
		}
		if meta.Message == "" {
			meta.Message = http.StatusText(resp.StatusCode())
		}
		return &Error{
			Response: errorResponse{
				Meta: Meta{
					Success: false,
					Message: "Invalid parameters or payload",
				},
				Error: E_UNPROCESSABLE_ENTITY,
			},
			Code:         http.StatusUnprocessableEntity,
			ErrorMessage: errors.New(meta.Message),
		}
	}
	return nil
}

func ErrorResponse(err error) *Error {
	var re *Error
	if errors.As(err, &re) {
		return re
	} else {
		return ErrorBuilder(&ErrorConstant.InternalServerError, err)
	}
}

func (e *Error) Error() string {
	if e.ErrorMessage == nil {
		e.ErrorMessage = errors.New(http.StatusText(e.Code))
	}
	return fmt.Sprintf("error code '%d' because: %s", e.Code, e.ErrorMessage.Error())
}

func (e *Error) ParseToError() error {
	return e
}

func (e *Error) WithData(data interface{}) *Error {
	e.Response.Error = data
	return e
}

func (e *Error) WithMetaMessage(message string) *Error {
	e.Response.Meta.Message = message
	return e
}

func (e *Error) Send(c echo.Context) error {
	if e.ErrorMessage != nil {
		logrus.Error(e.ErrorMessage)
	}
	if e.Header != nil {
		for k, values := range *e.Header {
			for _, v := range values {
				c.Response().Header().Add(k, v)
			}
		}
	}
	return c.JSON(e.Code, e.Response)
}
