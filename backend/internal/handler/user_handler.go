package handler

import (
	"fmt"
	"net/http"
	"time"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"
	"psychology-backend/pkg/util"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	*BaseHandler
	userService interfaces.UserServiceInterface
}

func NewUserHandler(userService interfaces.UserServiceInterface) *UserHandler {
	return &UserHandler{
		BaseHandler: NewBaseHandler(),
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var req schemas.UserCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	user, err := h.userService.CreateUser(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد کاربر", err.Error())
	}

	return response.Created(c, "کاربر با موفقیت ایجاد شد", user)
}

func (h *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه کاربر الزامی است")
	}

	user, err := h.userService.GetUserByID(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "کاربر")
	}

	return response.OK(c, "کاربر با موفقیت دریافت شد", user)
}

func (h *UserHandler) ListUsers(c echo.Context) error {
	var req schemas.UserListRequest
	if err := c.Bind(&req); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	users, err := h.userService.ListUsers(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت لیست کاربران", err)
	}

	return response.OK(c, "لیست کاربران با موفقیت دریافت شد", users)
}

func (h *UserHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه کاربر الزامی است")
	}

	var req schemas.UserUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	user, err := h.userService.UpdateUser(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی کاربر", err.Error())
	}

	return response.OK(c, "کاربر با موفقیت بروزرسانی شد", user)
}

func (h *UserHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه کاربر الزامی است")
	}

	if err := h.userService.DeleteUser(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف کاربر", err.Error())
	}

	return response.OK(c, "کاربر با موفقیت حذف شد", nil)
}

func (h *UserHandler) AcceptAgreement(c echo.Context) error {
	userID := h.GetUserID(c)
	if userID == "" {
		return h.Unauthorized(c, "شناسه کاربر یافت نشد")
	}

	user, err := h.userService.AcceptAgreement(c.Request().Context(), userID)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در پذیرش توافقنامه", err.Error())
	}

	return response.OK(c, "توافقنامه با موفقیت پذیرفته شد", user)
}

func (h *UserHandler) GetUserStats(c echo.Context) error {
	stats, err := h.userService.GetUserStats(c.Request().Context())
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار کاربران", err)
	}

	return response.OK(c, "آمار کاربران با موفقیت دریافت شد", stats)
}

func (h *UserHandler) GetUserActivityTrend(c echo.Context) error {
	var req struct {
		DateFrom string `query:"date_from" validate:"required"`
		DateTo   string `query:"date_to" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
	}

	dateFrom, err := parseDate(req.DateFrom)
	if err != nil {
		return h.BadRequest(c, "فرمت تاریخ شروع نامعتبر است")
	}

	dateTo, err := parseDate(req.DateTo)
	if err != nil {
		return h.BadRequest(c, "فرمت تاریخ پایان نامعتبر است")
	}

	trend, err := h.userService.GetUserActivityTrend(c.Request().Context(), dateFrom, dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت روند فعالیت کاربران", err)
	}

	return response.OK(c, "روند فعالیت کاربران با موفقیت دریافت شد", trend)
}

func (h *UserHandler) GetLoginAnalytics(c echo.Context) error {
	var req struct {
		DateFrom string `query:"date_from" validate:"required"`
		DateTo   string `query:"date_to" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
	}

	dateFrom, err := parseDate(req.DateFrom)
	if err != nil {
		return h.BadRequest(c, "فرمت تاریخ شروع نامعتبر است")
	}

	dateTo, err := parseDate(req.DateTo)
	if err != nil {
		return h.BadRequest(c, "فرمت تاریخ پایان نامعتبر است")
	}

	analytics, err := h.userService.GetLoginAnalytics(c.Request().Context(), dateFrom, dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار ورود", err)
	}

	return response.OK(c, "آمار ورود با موفقیت دریافت شد", analytics)
}

func (h *UserHandler) GetAgreementStats(c echo.Context) error {
	totalUsers, agreedUsers, rate, err := h.userService.GetAgreementStats(c.Request().Context())
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار توافقنامه", err)
	}

	return response.OK(c, "آمار توافقنامه با موفقیت دریافت شد", map[string]interface{}{
		"total_users":    totalUsers,
		"agreed_users":   agreedUsers,
		"agreement_rate": rate,
	})
}

func (h *UserHandler) GetAppVersionDistribution(c echo.Context) error {
	distribution, err := h.userService.GetAppVersionDistribution(c.Request().Context())
	if err != nil {
		return h.InternalError(c, "خطا در دریافت توزیع نسخه‌های اپلیکیشن", err)
	}

	return response.OK(c, "توزیع نسخه‌های اپلیکیشن با موفقیت دریافت شد", distribution)
}

func (h *UserHandler) GetInactiveUsers(c echo.Context) error {
	days := 30
	if daysStr := c.QueryParam("days"); daysStr != "" {
		if _, err := fmt.Sscanf(daysStr, "%d", &days); err != nil {
			return h.BadRequest(c, "فرمت تعداد روز نامعتبر است")
		}
	}

	users, err := h.userService.GetInactiveUsers(c.Request().Context(), days)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت کاربران غیرفعال", err)
	}

	return response.OK(c, "کاربران غیرفعال با موفقیت دریافت شدند", users)
}

func (h *UserHandler) GetUserEngagement(c echo.Context) error {
	var req struct {
		DateFrom string `query:"date_from" validate:"required"`
		DateTo   string `query:"date_to" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
	}

	dateFrom, err := parseDate(req.DateFrom)
	if err != nil {
		return h.BadRequest(c, "فرمت تاریخ شروع نامعتبر است")
	}

	dateTo, err := parseDate(req.DateTo)
	if err != nil {
		return h.BadRequest(c, "فرمت تاریخ پایان نامعتبر است")
	}

	stats, err := h.userService.GetUserEngagementStats(c.Request().Context(), dateFrom, dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار مشارکت", err)
	}

	return response.OK(c, "آمار مشارکت با موفقیت دریافت شد", stats)
}

func (h *UserHandler) ExportUsers(c echo.Context) error {
	var req struct {
		DateFrom *string `query:"date_from"`
		DateTo   *string `query:"date_to"`
		UserID   string  `query:"user_id"`
	}

	if err := c.Bind(&req); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
	}

	var dateFrom, dateTo *time.Time
	if req.DateFrom != nil {
		df, err := parseDate(*req.DateFrom)
		if err != nil {
			return h.BadRequest(c, "فرمت تاریخ شروع نامعتبر است")
		}
		dateFrom = &df
	}

	if req.DateTo != nil {
		dt, err := parseDate(*req.DateTo)
		if err != nil {
			return h.BadRequest(c, "فرمت تاریخ پایان نامعتبر است")
		}
		dateTo = &dt
	}

	users, err := h.userService.ExportUsers(c.Request().Context(), dateFrom, dateTo, req.UserID)
	if err != nil {
		return h.InternalError(c, "خطا در خروجی کاربران", err)
	}

	return response.OK(c, "خروجی کاربران با موفقیت دریافت شد", users)
}

func (h *UserHandler) GetUserByPhoneNumber(c echo.Context) error {
	phoneNumber := c.QueryParam("phone_number")
	if phoneNumber == "" {
		return h.BadRequest(c, "شماره تماس الزامی است")
	}

	user, err := h.userService.GetUserByPhoneNumber(c.Request().Context(), phoneNumber)
	if err != nil {
		return h.NotFound(c, "کاربر")
	}

	return response.OK(c, "کاربر با موفقیت دریافت شد", user)
}

func (h *UserHandler) UpdateLoginInfo(c echo.Context) error {
	userID := h.GetUserID(c)
	if userID == "" {
		return h.Unauthorized(c, "شناسه کاربر یافت نشد")
	}

	var req struct {
		AndroidVersion string `json:"android_version,omitempty"`
		AppVersion     string `json:"app_version,omitempty"`
	}

	if err := c.Bind(&req); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
	}

	user, err := h.userService.UpdateLoginInfo(c.Request().Context(), userID, req.AndroidVersion, req.AppVersion)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی اطلاعات ورود", err.Error())
	}

	return response.OK(c, "اطلاعات ورود با موفقیت بروزرسانی شد", user)
}

func parseDate(dateStr string) (time.Time, error) {
	return util.ParseShamsi(dateStr)
}

func (h *UserHandler) BadRequest(c echo.Context, message string) error {
	return response.BadRequest(c, message, "")
}

func (h *UserHandler) Unauthorized(c echo.Context, message string) error {
	return response.Unauthorized(c, message, "")
}

func (h *UserHandler) GetUserProfile(c echo.Context) error {
	userID := h.GetUserID(c)
	if userID == "" {
		return h.Unauthorized(c, "شناسه کاربر یافت نشد")
	}

	user, err := h.userService.GetUserProfile(c.Request().Context(), userID)
	if err != nil {
		return h.NotFound(c, "کاربر")
	}

	return response.OK(c, "پروفایل کاربر با موفقیت دریافت شد", user)
}

func (h *UserHandler) UpdateUserProfile(c echo.Context) error {
	userID := h.GetUserID(c)
	if userID == "" {
		return h.Unauthorized(c, "شناسه کاربر یافت نشد")
	}

	var req schemas.UserUpdateProfileRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	user, err := h.userService.UpdateUserProfile(c.Request().Context(), userID, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی پروفایل", err.Error())
	}

	return response.OK(c, "پروفایل کاربر با موفقیت بروزرسانی شد", user)
}

func (h *UserHandler) SyncUserData(c echo.Context) error {
	userID := h.GetUserID(c)
	if userID == "" {
		return h.Unauthorized(c, "شناسه کاربر یافت نشد")
	}

	var syncData map[string]interface{}
	if err := c.Bind(&syncData); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
	}

	if err := h.userService.SyncUserData(c.Request().Context(), userID, syncData); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در همگام‌سازی داده‌ها", err.Error())
	}

	return response.OK(c, "همگام‌سازی داده‌ها با موفقیت انجام شد", nil)
}
