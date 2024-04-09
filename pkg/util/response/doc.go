package response

// ErrorResponse400 ...
type ErrorResponse400 struct {
	Meta struct {
		Success bool   `json:"success" example:"false"`
		Message string `json:"message" example:"bad_request"`
	} `json:"meta"`
	Error interface{} `json:"data"`
}

// ErrorResponse401 ...
type ErrorResponse401 struct {
	Meta struct {
		Success bool   `json:"success" example:"false"`
		Message string `json:"message" example:"unauthorized"`
	} `json:"meta"`
	Error interface{} `json:"data"`
}

// ErrorResponse404 ...
type ErrorResponse404 struct {
	Meta struct {
		Success bool   `json:"success" example:"false"`
		Message string `json:"message" example:"not_found"`
	} `json:"meta"`
	Error interface{} `json:"data"`
}

// ErrorResponse422 ...
type ErrorResponse422 struct {
	Meta struct {
		Success bool   `json:"success" example:"false"`
		Message string `json:"message" example:"Invalid parameters or payload"`
	} `json:"meta"`
	Error string `json:"data" example:"unprocessable_entity"`
}

// ErrorResponse500 ...
type ErrorResponse500 struct {
	Meta struct {
		Success bool   `json:"success" example:"false"`
		Message string `json:"message" example:"Something bad happened"`
	} `json:"meta"`
	Error string `json:"data" example:"server_error"`
}
