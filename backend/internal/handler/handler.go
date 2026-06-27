package handler

import (
	"net/http"

	"psychology-backend/pkg/response"

	"github.com/labstack/echo/v4"
)

type BaseHandler struct{}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

func (h *BaseHandler) GetUserID(c echo.Context) string {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return ""
	}
	return userID
}

func (h *BaseHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}

func (h *BaseHandler) BindAndValidate(c echo.Context, req interface{}) error {
	if err := c.Bind(req); err != nil {
		return response.BadRequest(c, "بدنه درخواست نامعتبر است", err.Error())
	}

	if err := c.Validate(req); err != nil {
		return response.BadRequest(c, "اعتبارسنجی ناموفق بود", err.Error())
	}

	return nil
}

func (h *BaseHandler) PaginatedResponse(c echo.Context, message string, data interface{}, total int64, page, pageSize, pages int) error {
	return response.OK(c, message, map[string]interface{}{
		"data":      data,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"pages":     pages,
	})
}

func (h *BaseHandler) SuccessResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	return response.Success(c, statusCode, message, data)
}

func (h *BaseHandler) ErrorResponse(c echo.Context, statusCode int, message string, err string) error {
	return response.Error(c, statusCode, message, err)
}

func (h *BaseHandler) NotFound(c echo.Context, resource string) error {
	return response.Error(c, http.StatusNotFound, resource+" یافت نشد", "")
}

func (h *BaseHandler) InternalError(c echo.Context, message string, err error) error {
	return response.InternalServerError(c, message, err.Error())
}
