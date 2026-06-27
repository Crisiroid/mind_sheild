package handler

import (
	"net/http"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService interfaces.AuthServiceInterface
}

func NewAuthHandler(authService interfaces.AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) UserRegister(c echo.Context) error {
	var req schemas.UserRegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "بدنه درخواست نامعتبر است", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, "اعتبارسنجی ناموفق بود", err.Error())
	}

	result, err := h.authService.UserRegister(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "ثبت‌نام ناموفق بود", err.Error())
	}

	return response.Created(c, "کاربر با موفقیت ثبت‌نام شد", result)
}

func (h *AuthHandler) UserLogin(c echo.Context) error {
	var req schemas.UserLoginRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "بدنه درخواست نامعتبر است", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, "اعتبارسنجی ناموفق بود", err.Error())
	}

	result, err := h.authService.UserLogin(c.Request().Context(), &req)
	if err != nil {
		return response.Unauthorized(c, "ورود ناموفق بود", err.Error())
	}

	return response.OK(c, "ورود موفقیت‌آمیز بود", result)
}

func (h *AuthHandler) UserRefreshToken(c echo.Context) error {
	var req schemas.UserRefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "بدنه درخواست نامعتبر است", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, "اعتبارسنجی ناموفق بود", err.Error())
	}

	result, err := h.authService.UserRefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return response.Unauthorized(c, "تازه‌سازی توکن ناموفق بود", err.Error())
	}

	return response.OK(c, "توکن با موفقیت تازه‌سازی شد", result)
}

func (h *AuthHandler) UserLogout(c echo.Context) error {
	userID := c.Get("user_id").(string)

	if err := h.authService.UserLogout(c.Request().Context(), userID); err != nil {
		return response.Error(c, http.StatusBadRequest, "خروج ناموفق بود", err.Error())
	}

	return response.OK(c, "با موفقیت خارج شدید", nil)
}

func (h *AuthHandler) UserChangePassword(c echo.Context) error {
	userID := c.Get("user_id").(string)

	var req schemas.UserChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "بدنه درخواست نامعتبر است", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, "اعتبارسنجی ناموفق بود", err.Error())
	}

	if err := h.authService.UserChangePassword(c.Request().Context(), userID, &req); err != nil {
		return response.Error(c, http.StatusBadRequest, "تغییر رمز عبور ناموفق بود", err.Error())
	}

	return response.OK(c, "رمز عبور با موفقیت تغییر کرد", nil)
}

func (h *AuthHandler) AdminLogin(c echo.Context) error {
	var req schemas.AdminLoginRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "بدنه درخواست نامعتبر است", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, "اعتبارسنجی ناموفق بود", err.Error())
	}

	result, err := h.authService.AdminLogin(c.Request().Context(), &req)
	if err != nil {
		return response.Unauthorized(c, "ورود ناموفق بود", err.Error())
	}

	return response.OK(c, "ورود موفقیت‌آمیز بود", result)
}

func (h *AuthHandler) AdminRefreshToken(c echo.Context) error {
	var req schemas.AdminRefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "بدنه درخواست نامعتبر است", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, "اعتبارسنجی ناموفق بود", err.Error())
	}

	result, err := h.authService.AdminRefreshToken(c.Request().Context(), req.RefreshToken)
	if err != nil {
		return response.Unauthorized(c, "تازه‌سازی توکن ناموفق بود", err.Error())
	}

	return response.OK(c, "توکن با موفقیت تازه‌سازی شد", result)
}

func (h *AuthHandler) AdminLogout(c echo.Context) error {
	adminID := c.Get("user_id").(string)

	if err := h.authService.AdminLogout(c.Request().Context(), adminID); err != nil {
		return response.Error(c, http.StatusBadRequest, "خروج ناموفق بود", err.Error())
	}

	return response.OK(c, "با موفقیت خارج شدید", nil)
}

func (h *AuthHandler) AdminChangePassword(c echo.Context) error {
	adminID := c.Get("user_id").(string)

	var req schemas.AdminChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "بدنه درخواست نامعتبر است", err.Error())
	}

	if err := c.Validate(&req); err != nil {
		return response.BadRequest(c, "اعتبارسنجی ناموفق بود", err.Error())
	}

	if err := h.authService.AdminChangePassword(c.Request().Context(), adminID, &req); err != nil {
		return response.Error(c, http.StatusBadRequest, "تغییر رمز عبور ناموفق بود", err.Error())
	}

	return response.OK(c, "رمز عبور با موفقیت تغییر کرد", nil)
}
