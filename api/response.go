package api

// ErrorResponse represents error response body
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}
