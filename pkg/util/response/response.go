package response

import "boilerplate/internal/abstraction"

type Meta struct {
	Success bool                        `json:"success" default:"true"`
	Message string                      `json:"message"`
	Info    *abstraction.PaginationInfo `json:"info"`
}
