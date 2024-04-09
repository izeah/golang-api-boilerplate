package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

func TestError_Send(t *testing.T) {
	tmr := ErrorConstant.TooManyRequest(10)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	type fields struct {
		Header       *http.Header
		Response     errorResponse
		Code         int
		ErrorMessage error
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "should success",
			fields: fields{
				Header:       tmr.Header,
				Response:     tmr.Response,
				Code:         tmr.Code,
				ErrorMessage: tmr.ErrorMessage,
			},
			args: args{
				c: echo.New().NewContext(req, rec),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Header:       tt.fields.Header,
				Response:     tt.fields.Response,
				Code:         tt.fields.Code,
				ErrorMessage: tt.fields.ErrorMessage,
			}
			if err := e.Send(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCustomErrorBuilder(t *testing.T) {
	type args struct {
		code    int
		err     interface{}
		message string
		vals    []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "Unauthorized",
			args: args{
				code:    http.StatusUnauthorized,
				err:     nil,
				message: "unauthorized",
				vals:    nil,
			},
			want: nil,
		},
		{
			name: "battery_tx_notfound",
			args: args{
				code:    http.StatusNotFound,
				err:     nil,
				message: "battery_tx_notfound",
				vals:    nil,
			},
			want: nil,
		},
		{
			name: "door_empty_notfound",
			args: args{
				code:    http.StatusNotFound,
				err:     nil,
				message: "door_empty_notfound",
				vals:    nil,
			},
			want: nil,
		},
		{
			name: "door_full_and_battery_notfound",
			args: args{
				code:    http.StatusNotFound,
				err:     nil,
				message: "door_full_and_battery_notfound",
				vals:    nil,
			},
			want: nil,
		},
		{
			name: "customer_package_usage_notfound",
			args: args{
				code:    http.StatusNotFound,
				err:     nil,
				message: "customer_package_usage_notfound",
				vals:    nil,
			},
			want: nil,
		},
		{
			name: "stb_notfound",
			args: args{
				code:    http.StatusNotFound,
				err:     nil,
				message: "stb_notfound",
				vals:    nil,
			},
			want: nil,
		},
		{
			name: "should success",
			args: args{
				code:    500,
				err:     nil,
				message: "something wrong",
				vals:    nil,
			},
			want: nil,
		},
		{
			name: "should success",
			args: args{
				code:    500,
				err:     errors.New("something wrong"),
				message: "something wrong",
				vals:    nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CustomErrorBuilder(tt.args.code, tt.args.err, tt.args.message, tt.args.vals...); !reflect.DeepEqual(got, tt.want) {
				b, _ := json.MarshalIndent(got.Response, "", " ")
				fmt.Println(string(b))
			}
		})
	}
}

func TestCustomErrorBuilder2(t *testing.T) {
	type args struct {
		code    int
		data    interface{}
		message string
		vals    []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "customermotorcycle_notfound",
			args: args{
				code:    http.StatusNotFound,
				data:    nil,
				message: "customermotorcycle_notfound",
				vals:    nil,
			},
			want: nil,
		}, {
			name: "cabinet_in_use",
			args: args{
				code: http.StatusNotFound,
				data: map[string]interface{}{
					"name":          "Volter",
					"duration":      1,
					"duration_unit": "minute",
				},
				message: "cabinet_in_use",
				vals:    nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CustomErrorBuilder(tt.args.code, tt.args.data, tt.args.message, tt.args.vals...); !reflect.DeepEqual(got, tt.want) {
				b, _ := json.MarshalIndent(got.Response, "", " ")
				fmt.Println(string(b))
			}
		})
	}
}

func TestErrorBuilder(t *testing.T) {
	type args struct {
		res     *Error
		message error
		vals    []interface{}
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		{
			name: "too_many_request",
			args: args{
				res:     ErrorConstant.TooManyRequest(1),
				message: fmt.Errorf("too many request, last trx %v - %v", 1, time.Now()),
				vals:    nil,
			},
			want: nil,
		},
		{
			name: "bad_request",
			args: args{
				res:     &ErrorConstant.BadRequest,
				message: errors.New("something wrong"),
				vals:    nil,
			},
			want: nil,
		},
		{
			name: "validation",
			args: args{
				res:     &ErrorConstant.Validation,
				message: errors.New("something wrong"),
				vals:    nil,
			},
			want: nil,
		},
		{
			name: "internal_server_error",
			args: args{
				res:     &ErrorConstant.InternalServerError,
				message: errors.New("something wrong"),
				vals:    nil,
			},
			want: nil,
		},
		{
			name: "unprocessable_entity",
			args: args{
				res:     &ErrorConstant.UnprocessableEntity,
				message: errors.New("something wrong"),
				vals:    nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ErrorBuilder(tt.args.res, tt.args.message, tt.args.vals...); !reflect.DeepEqual(got, tt.want) {
				if tt.name == "too_many_request" {
					got = got.WithData(map[string]interface{}{
						"duration":      1,
						"duration_unit": "second",
					}).WithMetaMessage("too_many_request")
				}

				b, _ := json.MarshalIndent(got.Response, "", " ")
				fmt.Println(string(b))
			}
		})
	}
}
