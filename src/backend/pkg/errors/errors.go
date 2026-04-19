package errors

import (
	"fmt"
	"net/http"
)

// Type represents the domain error type.
type Type string

const (
	TypeNotFound     Type = "NOT_FOUND"
	TypeValidation   Type = "VALIDATION"
	TypeConflict     Type = "CONFLICT"
	TypeUnauthorized Type = "UNAUTHORIZED"
	TypeInternal     Type = "INTERNAL_ERROR"
)

// AppError is the standard application error.
type AppError struct {
	Code    Type                   `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	err     error
}

func (e *AppError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.err
}

// New creates a new AppError.
func New(code Type, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		err:     err,
	}
}

// ErrorResponse is the HTTP response structure for errors.
type ErrorResponse struct {
	Error AppErrorResponse `json:"error"`
}

type AppErrorResponse struct {
	Code    Type                   `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// ToHTTPStatus maps the domain error type to HTTP status codes.
func ToHTTPStatus(code Type) int {
	switch code {
	case TypeNotFound:
		return http.StatusNotFound
	case TypeValidation:
		return http.StatusBadRequest
	case TypeConflict:
		return http.StatusConflict
	case TypeUnauthorized:
		return http.StatusUnauthorized
	case TypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Handle maps a Go error to an ErrorResponse and HTTP status code.
func Handle(err error, env string) (int, ErrorResponse) {
	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	} else {
		appErr = New(TypeInternal, "Internal server error", err)
	}

	resp := ErrorResponse{
		Error: AppErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
	}

	// Never expose stack trace or internal details in production
	if appErr.Code == TypeInternal && env == "production" {
		resp.Error.Message = "Internal server error"
		resp.Error.Details = nil
	} else if appErr.Code == TypeInternal && appErr.err != nil {
		if resp.Error.Details == nil {
			resp.Error.Details = make(map[string]interface{})
		}
		resp.Error.Details["internal_error"] = appErr.err.Error()
	}

	return ToHTTPStatus(appErr.Code), resp
}
