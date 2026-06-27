package handler

import (
	"net/http"
	"strconv"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type CalendarHandler struct {
	*BaseHandler
	calendarService interfaces.DailyCalendarServiceInterface
}

func NewCalendarHandler(calendarService interfaces.DailyCalendarServiceInterface) *CalendarHandler {
	return &CalendarHandler{
		BaseHandler:     NewBaseHandler(),
		calendarService: calendarService,
	}
}

func (h *CalendarHandler) CreateCalendarEntry(c echo.Context) error {
	var req schemas.CalendarCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	calendar, err := h.calendarService.CreateCalendarEntry(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد تقویم روزانه", err.Error())
	}

	return response.Created(c, "تقویم روزانه با موفقیت ایجاد شد", calendar)
}

func (h *CalendarHandler) GetCalendarEntryByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه تقویم الزامی است")
	}

	calendar, err := h.calendarService.GetCalendarEntryById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "تقویم روزانه")
	}

	return response.OK(c, "تقویم روزانه با موفقیت دریافت شد", calendar)
}

func (h *CalendarHandler) ListCalendarEntries(c echo.Context) error {
	var req schemas.CalendarListRequest
	if err := c.Bind(&req); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
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

	calendars, err := h.calendarService.GetAllCalendarEntries(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت تقویم‌های روزانه", err)
	}

	return response.OK(c, "تقویم‌های روزانه با موفقیت دریافت شدند", calendars)
}

func (h *CalendarHandler) UpdateCalendarEntry(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه تقویم الزامی است")
	}

	var req schemas.CalendarUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	calendar, err := h.calendarService.UpdateCalendarEntry(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی تقویم روزانه", err.Error())
	}

	return response.OK(c, "تقویم روزانه با موفقیت بروزرسانی شد", calendar)
}

func (h *CalendarHandler) DeleteCalendarEntry(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه تقویم الزامی است")
	}

	if err := h.calendarService.DeleteCalendarEntry(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف تقویم روزانه", err.Error())
	}

	return response.OK(c, "تقویم روزانه با موفقیت حذف شد", nil)
}

func (h *CalendarHandler) GetCompletionStats(c echo.Context) error {
	userID := h.GetUserID(c)
	if userID == "" {
		return h.Unauthorized(c, "شناسه کاربر یافت نشد")
	}

	stats, err := h.calendarService.GetCompletionStats(c.Request().Context(), userID)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار تکمیل", err)
	}

	return response.OK(c, "آمار تکمیل با موفقیت دریافت شد", stats)
}

func (h *CalendarHandler) GetDayRangeProgress(c echo.Context) error {
	userID := h.GetUserID(c)
	if userID == "" {
		return h.Unauthorized(c, "شناسه کاربر یافت نشد")
	}

	fromDayStr := c.QueryParam("from_day")
	toDayStr := c.QueryParam("to_day")

	if fromDayStr == "" || toDayStr == "" {
		return h.BadRequest(c, "پارامترهای from_day و to_day الزامی هستند")
	}

	fromDay, err := strconv.Atoi(fromDayStr)
	if err != nil || fromDay < 1 || fromDay > 56 {
		return h.BadRequest(c, "فرمت from_day نامعتبر است (1-56)")
	}

	toDay, err := strconv.Atoi(toDayStr)
	if err != nil || toDay < 1 || toDay > 56 {
		return h.BadRequest(c, "فرمت to_day نامعتبر است (1-56)")
	}

	if fromDay > toDay {
		return h.BadRequest(c, "from_day باید کوچکتر از to_day باشد")
	}

	progress, err := h.calendarService.GetDayRangeProgress(c.Request().Context(), userID, fromDay, toDay)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت پیشرفت", err)
	}

	return response.OK(c, "پیشرفت با موفقیت دریافت شد", progress)
}

func (h *CalendarHandler) GetStreakAnalysis(c echo.Context) error {
	userID := h.GetUserID(c)
	if userID == "" {
		return h.Unauthorized(c, "شناسه کاربر یافت نشد")
	}

	analysis, err := h.calendarService.GetStreakAnalysis(c.Request().Context(), userID)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت تحلیل پشت سر هم", err)
	}

	return response.OK(c, "تحلیل پشت سر هم با موفقیت دریافت شد", analysis)
}

func (h *CalendarHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}

func (h *CalendarHandler) BadRequest(c echo.Context, message string) error {
	return response.BadRequest(c, message, "")
}

func (h *CalendarHandler) Unauthorized(c echo.Context, message string) error {
	return response.Unauthorized(c, message, "")
}
