package utils

import (
	"fmt"
	"time"
)

type CustomError struct {
	Code    string
	Message string
	Details interface{}
}

func (e *CustomError) Error() string {
	return e.Message
}

// Common error constructors
func NewValidationError(details interface{}) *CustomError {
	return &CustomError{
		Code:    "VALIDATION_ERROR",
		Message: "Request validation failed",
		Details: details,
	}
}

func NewUnauthorizedError() *CustomError {
	return &CustomError{
		Code:    "UNAUTHORIZED",
		Message: "Unauthorized access",
	}
}

func NewForbiddenError() *CustomError {
	return &CustomError{
		Code:    "FORBIDDEN",
		Message: "Access forbidden",
	}
}

func NewNotFoundError(resource string) *CustomError {
	return &CustomError{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s not found", resource),
	}
}

func NewInternalError(message string) *CustomError {
	return &CustomError{
		Code:    "INTERNAL_ERROR",
		Message: message,
	}
}

func NewConflictError(message string) *CustomError {
	return &CustomError{
		Code:    "CONFLICT",
		Message: message,
	}
}

func NewCustomError(code, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

// Response wrapper
type Response struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp string      `json:"timestamp"`
}

func NewSuccessResponse(data interface{}, message string) Response {
	return Response{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func NewErrorResponse(err *CustomError) Response {
	return Response{
		Success: false,
		Error: map[string]interface{}{
			"code":    err.Code,
			"message": err.Message,
			"details": err.Details,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
