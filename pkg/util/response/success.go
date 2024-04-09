package response

import (
	"net/http"

	"boilerplate/internal/abstraction"

	"github.com/labstack/echo/v4"
)

type successConstant struct {
	OK Success
}

var SuccessConstant = successConstant{
	OK: Success{
		Response: successResponse{
			Meta: Meta{
				Success: true,
				Message: "Request successfully proceed",
			},
			Data: nil,
		},
		Code: http.StatusOK,
	},
}

type successResponse struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Success struct {
	Response successResponse `json:"response"`
	Code     int             `json:"code"`
}

func SuccessBuilder(res *Success, data interface{}) *Success {
	res.Response.Data = data
	return res
}

func SuccessResponse(data interface{}) *Success {
	return SuccessBuilder(&SuccessConstant.OK, data)
}

func (s *Success) WithPagination(info *abstraction.PaginationInfo) *Success {
	s.Response.Meta.Info = info
	return s
}

func (s *Success) WithRowsAffected(rowsAffected int64) *Success {
	code := http.StatusOK
	if rowsAffected == 1 {
		code = http.StatusCreated
	}
	s.Code = code
	return s
}

func (s *Success) Send(c echo.Context) error {
	return c.JSON(s.Code, s.Response)
}
