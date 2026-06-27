package handler

import (
	"net/http"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type RoleValueHandler struct {
	*BaseHandler
	roleValueService interfaces.RoleValueServiceInterface
}

func NewRoleValueHandler(roleValueService interfaces.RoleValueServiceInterface) *RoleValueHandler {
	return &RoleValueHandler{
		BaseHandler:      NewBaseHandler(),
		roleValueService: roleValueService,
	}
}

func (h *RoleValueHandler) CreateRoleValue(c echo.Context) error {
	var req schemas.RoleValueCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	roleValue, err := h.roleValueService.CreateRoleValue(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد نقش/ارزش", err.Error())
	}

	return response.Created(c, "نقش/ارزش با موفقیت ایجاد شد", roleValue)
}

func (h *RoleValueHandler) GetRoleValueByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه نقش/ارزش الزامی است", "")
	}

	roleValue, err := h.roleValueService.GetRoleValueById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "نقش/ارزش")
	}

	return response.OK(c, "نقش/ارزش با موفقیت دریافت شد", roleValue)
}

func (h *RoleValueHandler) ListRolesValues(c echo.Context) error {
	var req schemas.RoleValueListRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "بدنه درخواست نامعتبر است", "")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	if h.GetUserRole(c) == "user" {
		req.UserID = h.GetUserID(c)
	}

	rolesValue, err := h.roleValueService.GetAllRoleValues(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت نقش‌ها/ارزش‌ها", err)
	}

	return response.OK(c, "نقش‌ها/ارزش‌ها با موفقیت دریافت شدند", rolesValue)
}

func (h *RoleValueHandler) UpdateRoleValue(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه نقش/ارزش الزامی است", "")
	}

	var req schemas.RoleValueUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	roleValue, err := h.roleValueService.UpdateRoleValue(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی نقش/ارزش", err.Error())
	}

	return response.OK(c, "نقش/ارزش با موفقیت بروزرسانی شد", roleValue)
}

func (h *RoleValueHandler) DeleteRoleValue(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه نقش/ارزش الزامی است", "")
	}

	if err := h.roleValueService.DeleteRoleValue(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف نقش/ارزش", err.Error())
	}

	return response.OK(c, "نقش/ارزش با موفقیت حذف شد", nil)
}

func (h *RoleValueHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}
