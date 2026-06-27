package handler

import (
	"net/http"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	*BaseHandler
	adminUserService interfaces.AdminUserServiceInterface
	adminRoleService interfaces.AdminRoleServiceInterface
	systemLogService interfaces.SystemLogServiceInterface
}

func NewAdminHandler(adminUserService interfaces.AdminUserServiceInterface, adminRoleService interfaces.AdminRoleServiceInterface, systemLogService interfaces.SystemLogServiceInterface) *AdminHandler {
	return &AdminHandler{
		BaseHandler:      NewBaseHandler(),
		adminUserService: adminUserService,
		adminRoleService: adminRoleService,
		systemLogService: systemLogService,
	}
}

func (h *AdminHandler) CreateAdminUser(c echo.Context) error {
	var req schemas.AdminCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	admin, err := h.adminUserService.CreateAdminUser(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد مدیر", err.Error())
	}

	return response.Created(c, "مدیر با موفقیت ایجاد شد", admin)
}

func (h *AdminHandler) GetAdminUserByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه مدیر الزامی است")
	}

	admin, err := h.adminUserService.GetAdminUserById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "مدیر")
	}

	return response.OK(c, "مدیر با موفقیت دریافت شد", admin)
}

func (h *AdminHandler) ListAdminUsers(c echo.Context) error {
	var req schemas.AdminListRequest
	if err := c.Bind(&req); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	admins, err := h.adminUserService.GetAllAdmins(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت مدیران", err)
	}

	return response.OK(c, "مدیران با موفقیت دریافت شدند", admins)
}

func (h *AdminHandler) UpdateAdminUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه مدیر الزامی است")
	}

	var req schemas.AdminUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	admin, err := h.adminUserService.UpdateAdminUser(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی مدیر", err.Error())
	}

	return response.OK(c, "مدیر با موفقیت بروزرسانی شد", admin)
}

func (h *AdminHandler) DeleteAdminUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه مدیر الزامی است")
	}

	if err := h.adminUserService.DeleteAdminUser(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف مدیر", err.Error())
	}

	return response.OK(c, "مدیر با موفقیت حذف شد", nil)
}

func (h *AdminHandler) DeactivateAdminUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه مدیر الزامی است")
	}

	admin, err := h.adminUserService.DeactivateAdminUser(c.Request().Context(), id)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در غیرفعال کردن مدیر", err.Error())
	}

	return response.OK(c, "مدیر با موفقیت غیرفعال شد", admin)
}

func (h *AdminHandler) CreateAdminRole(c echo.Context) error {
	var req schemas.AdminRoleCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	role, err := h.adminRoleService.CreateAdminRole(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد نقش", err.Error())
	}

	return response.Created(c, "نقش با موفقیت ایجاد شد", role)
}

func (h *AdminHandler) GetAdminRoleByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه نقش الزامی است")
	}

	role, err := h.adminRoleService.GetAdminRoleById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "نقش")
	}

	return response.OK(c, "نقش با موفقیت دریافت شد", role)
}

func (h *AdminHandler) ListAdminRoles(c echo.Context) error {
	var req schemas.AdminRoleListRequest
	if err := c.Bind(&req); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	roles, err := h.adminRoleService.GetAllRoles(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت نقش‌ها", err)
	}

	return response.OK(c, "نقش‌ها با موفقیت دریافت شدند", roles)
}

func (h *AdminHandler) UpdateAdminRole(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه نقش الزامی است")
	}

	var req schemas.AdminRoleUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	role, err := h.adminRoleService.UpdateAdminRole(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی نقش", err.Error())
	}

	return response.OK(c, "نقش با موفقیت بروزرسانی شد", role)
}

func (h *AdminHandler) DeleteAdminRole(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه نقش الزامی است")
	}

	if err := h.adminRoleService.DeleteAdminRole(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف نقش", err.Error())
	}

	return response.OK(c, "نقش با موفقیت حذف شد", nil)
}

func (h *AdminHandler) ListSystemLogs(c echo.Context) error {
	var req schemas.SystemLogListRequest
	if err := c.Bind(&req); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	logs, err := h.systemLogService.GetAllSystemLogs(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت لاگ‌های سیستم", err)
	}

	return response.OK(c, "لاگ‌های سیستم با موفقیت دریافت شدند", logs)
}

func (h *AdminHandler) GetSystemLogByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه لاگ الزامی است")
	}

	log, err := h.systemLogService.GetSystemLogById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "لاگ سیستم")
	}

	return response.OK(c, "لاگ سیستم با موفقیت دریافت شد", log)
}

func (h *AdminHandler) BadRequest(c echo.Context, message string) error {
	return response.BadRequest(c, message, "")
}

func (h *AdminHandler) Unauthorized(c echo.Context, message string) error {
	return response.Unauthorized(c, message, "")
}

func (h *AdminHandler) NotFound(c echo.Context, resource string) error {
	return response.Error(c, http.StatusNotFound, resource+" یافت نشد", "")
}

func (h *AdminHandler) InternalError(c echo.Context, message string, err error) error {
	return response.InternalServerError(c, message, err.Error())
}

func (h *AdminHandler) GetAdminProfile(c echo.Context) error {
	adminID := h.GetUserID(c)
	if adminID == "" {
		return h.Unauthorized(c, "شناسه مدیر یافت نشد")
	}

	admin, err := h.adminUserService.GetAdminProfile(c.Request().Context(), adminID)
	if err != nil {
		return h.NotFound(c, "مدیر")
	}

	return response.OK(c, "پروفایل مدیر با موفقیت دریافت شد", admin)
}

func (h *AdminHandler) UpdateAdminProfile(c echo.Context) error {
	adminID := h.GetUserID(c)
	if adminID == "" {
		return h.Unauthorized(c, "شناسه مدیر یافت نشد")
	}

	var req schemas.AdminUpdateProfileRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	admin, err := h.adminUserService.UpdateAdminProfile(c.Request().Context(), adminID, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی پروفایل مدیر", err.Error())
	}

	return response.OK(c, "پروفایل مدیر با موفقیت بروزرسانی شد", admin)
}
