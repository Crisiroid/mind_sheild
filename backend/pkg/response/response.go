package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func Success(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c echo.Context, statusCode int, message string, err string) error {
	return c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func Created(c echo.Context, message string, data interface{}) error {
	return Success(c, http.StatusCreated, message, data)
}

func OK(c echo.Context, message string, data interface{}) error {
	return Success(c, http.StatusOK, message, data)
}

func BadRequest(c echo.Context, message string, err string) error {
	return Error(c, http.StatusBadRequest, message, err)
}

func Unauthorized(c echo.Context, message string, err string) error {
	return Error(c, http.StatusUnauthorized, message, err)
}

func InternalServerError(c echo.Context, message string, err string) error {
	return Error(c, http.StatusInternalServerError, message, err)
}
