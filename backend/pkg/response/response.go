package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Success sends a success response
func Success(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c echo.Context, statusCode int, message string, err string) error {
	return c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

// Created sends a created response
func Created(c echo.Context, message string, data interface{}) error {
	return Success(c, http.StatusCreated, message, data)
}

// OK sends an OK response
func OK(c echo.Context, message string, data interface{}) error {
	return Success(c, http.StatusOK, message, data)
}

// BadRequest sends a bad request response
func BadRequest(c echo.Context, message string, err string) error {
	return Error(c, http.StatusBadRequest, message, err)
}

// Unauthorized sends an unauthorized response
func Unauthorized(c echo.Context, message string, err string) error {
	return Error(c, http.StatusUnauthorized, message, err)
}

// InternalServerError sends an internal server error response
func InternalServerError(c echo.Context, message string, err string) error {
	return Error(c, http.StatusInternalServerError, message, err)
}
