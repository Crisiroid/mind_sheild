package handler

import (
	"net/http"
	"time"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type ReportHandler struct {
	*BaseHandler
	userReportService    interfaces.UserReportServiceInterface
	weeklyReportService  interfaces.WeeklyReportServiceInterface
	userService          interfaces.UserServiceInterface
	stressEventService   interfaces.StressEventServiceInterface
	bodyTensionService   interfaces.BodyTensionServiceInterface
	cognitiveGameService interfaces.CognitiveGameServiceInterface
	moodTrackerService   interfaces.MoodTrackerServiceInterface
	breathingService     interfaces.BreathingServiceInterface
}

func NewReportHandler(
	userReportService interfaces.UserReportServiceInterface,
	weeklyReportService interfaces.WeeklyReportServiceInterface,
	userService interfaces.UserServiceInterface,
	stressEventService interfaces.StressEventServiceInterface,
	bodyTensionService interfaces.BodyTensionServiceInterface,
	cognitiveGameService interfaces.CognitiveGameServiceInterface,
	moodTrackerService interfaces.MoodTrackerServiceInterface,
	breathingService interfaces.BreathingServiceInterface,
) *ReportHandler {
	return &ReportHandler{
		BaseHandler:          NewBaseHandler(),
		userReportService:    userReportService,
		weeklyReportService:  weeklyReportService,
		userService:          userService,
		stressEventService:   stressEventService,
		bodyTensionService:   bodyTensionService,
		cognitiveGameService: cognitiveGameService,
		moodTrackerService:   moodTrackerService,
		breathingService:     breathingService,
	}
}

func (h *ReportHandler) CreateUserReport(c echo.Context) error {
	var req schemas.UserReportCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	report, err := h.userReportService.CreateUserReport(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد گزارش کاربر", err.Error())
	}

	return response.Created(c, "گزارش کاربر با موفقیت ایجاد شد", report)
}

func (h *ReportHandler) GetUserReportByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه گزارش الزامی است")
	}

	report, err := h.userReportService.GetUserReportById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "گزارش کاربر")
	}

	return response.OK(c, "گزارش کاربر با موفقیت دریافت شد", report)
}

func (h *ReportHandler) ListUserReports(c echo.Context) error {
	var req schemas.UserReportListRequest
	if err := c.Bind(&req); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	reports, err := h.userReportService.GetAllUserReports(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت گزارش‌های کاربر", err)
	}

	return response.OK(c, "گزارش‌های کاربر با موفقیت دریافت شدند", reports)
}

func (h *ReportHandler) DeleteUserReport(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه گزارش الزامی است")
	}

	if err := h.userReportService.DeleteUserReport(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف گزارش کاربر", err.Error())
	}

	return response.OK(c, "گزارش کاربر با موفقیت حذف شد", nil)
}

func (h *ReportHandler) CreateWeeklyReport(c echo.Context) error {
	var req schemas.WeeklyReportCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	report, err := h.weeklyReportService.CreateWeeklyReport(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد گزارش هفتگی", err.Error())
	}

	return response.Created(c, "گزارش هفتگی با موفقیت ایجاد شد", report)
}

func (h *ReportHandler) GetWeeklyReportByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه گزارش الزامی است")
	}

	report, err := h.weeklyReportService.GetWeeklyReportById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "گزارش هفتگی")
	}

	return response.OK(c, "گزارش هفتگی با موفقیت دریافت شد", report)
}

func (h *ReportHandler) ListWeeklyReports(c echo.Context) error {
	var req schemas.WeeklyReportListRequest
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

	reports, err := h.weeklyReportService.GetAllWeeklyReports(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت گزارش‌های هفتگی", err)
	}

	return response.OK(c, "گزارش‌های هفتگی با موفقیت دریافت شدند", reports)
}

func (h *ReportHandler) UpdateWeeklyReport(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه گزارش الزامی است")
	}

	var req schemas.WeeklyReportUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	report, err := h.weeklyReportService.UpdateWeeklyReport(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی گزارش هفتگی", err.Error())
	}

	return response.OK(c, "گزارش هفتگی با موفقیت بروزرسانی شد", report)
}

func (h *ReportHandler) DeleteWeeklyReport(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه گزارش الزامی است")
	}

	if err := h.weeklyReportService.DeleteWeeklyReport(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف گزارش هفتگی", err.Error())
	}

	return response.OK(c, "گزارش هفتگی با موفقیت حذف شد", nil)
}

func (h *ReportHandler) GetDashboard(c echo.Context) error {
	userID := h.GetUserID(c)
	role := h.GetUserRole(c)

	if role == "user" {
		return h.getUserDashboard(c, userID)
	}

	return h.getAdminDashboard(c)
}

func (h *ReportHandler) getUserDashboard(c echo.Context, userID string) error {
	stats, err := h.userService.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت داشبورد", err)
	}

	return response.OK(c, "داشبورد با موفقیت دریافت شد", stats)
}

func (h *ReportHandler) getAdminDashboard(c echo.Context) error {
	userStats, err := h.userService.GetUserStats(c.Request().Context())
	if err != nil {
		return h.InternalError(c, "خطا در دریافت داشبورد", err)
	}

	return response.OK(c, "داشبورد با موفقیت دریافت شد", map[string]interface{}{
		"user_stats": userStats,
	})
}

func (h *ReportHandler) GetUserActivity(c echo.Context) error {
	ctx := c.Request().Context()
	now := time.Now()
	dateFrom := now.AddDate(0, 0, -30)
	dateTo := now

	if c.QueryParam("date_from") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_from")); err == nil {
			dateFrom = parsed
		}
	}
	if c.QueryParam("date_to") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_to")); err == nil {
			dateTo = parsed
		}
	}

	registrationTrend, err := h.userService.GetUserActivityTrend(ctx, dateFrom, dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت روند فعالیت کاربران", err)
	}

	loginTrend, err := h.userService.GetLoginAnalytics(ctx, dateFrom, dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار ورود", err)
	}

	breathingStats, err := h.breathingService.GetSessionStats(ctx, "", &dateFrom, &dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار تنفس", err)
	}

	return response.OK(c, "فعالیت کاربران دریافت شد", map[string]interface{}{
		"date_range":         map[string]string{"from": dateFrom.Format("2006-01-02"), "to": dateTo.Format("2006-01-02")},
		"registration_trend": registrationTrend,
		"login_trend":        loginTrend,
		"breathing_sessions": breathingStats,
	})
}

func (h *ReportHandler) GetStressAnalytics(c echo.Context) error {
	ctx := c.Request().Context()
	now := time.Now()
	dateFrom := now.AddDate(0, 0, -30)
	dateTo := now

	if c.QueryParam("date_from") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_from")); err == nil {
			dateFrom = parsed
		}
	}
	if c.QueryParam("date_to") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_to")); err == nil {
			dateTo = parsed
		}
	}

	stats, err := h.stressEventService.GetStressStats(ctx, "", &dateFrom, &dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در تحلیل استرس", err)
	}

	trend, err := h.stressEventService.GetIntensityTrend(ctx, "", dateFrom, dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت روند استرس", err)
	}

	distribution, err := h.stressEventService.GetSituationTypeDistribution(ctx, "")
	if err != nil {
		return h.InternalError(c, "خطا در دریافت توزیع موقعیت‌ها", err)
	}

	return response.OK(c, "تحلیل استرس دریافت شد", map[string]interface{}{
		"date_range":             map[string]string{"from": dateFrom.Format("2006-01-02"), "to": dateTo.Format("2006-01-02")},
		"stats":                  stats,
		"intensity_trend":        trend,
		"situation_distribution": distribution,
	})
}

func (h *ReportHandler) GetBodyTensionReport(c echo.Context) error {
	ctx := c.Request().Context()
	now := time.Now()
	dateFrom := now.AddDate(0, 0, -30)
	dateTo := now

	if c.QueryParam("date_from") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_from")); err == nil {
			dateFrom = parsed
		}
	}
	if c.QueryParam("date_to") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_to")); err == nil {
			dateTo = parsed
		}
	}

	stats, err := h.bodyTensionService.GetAverageIntensity(ctx, "", &dateFrom, &dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار تنش بدنی", err)
	}

	trend, err := h.bodyTensionService.GetIntensityTrend(ctx, "", dateFrom, dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت روند تنش بدنی", err)
	}

	colorDistribution, err := h.bodyTensionService.GetSeverityColorDistribution(ctx, "")
	if err != nil {
		return h.InternalError(c, "خطا در دریافت توزیع رنگ‌ها", err)
	}

	return response.OK(c, "گزارش تنش بدنی دریافت شد", map[string]interface{}{
		"date_range":            map[string]string{"from": dateFrom.Format("2006-01-02"), "to": dateTo.Format("2006-01-02")},
		"stats":                 stats,
		"intensity_trend":       trend,
		"severity_distribution": colorDistribution,
	})
}

func (h *ReportHandler) GetCognitivePatterns(c echo.Context) error {
	ctx := c.Request().Context()
	now := time.Now()
	dateFrom := now.AddDate(0, 0, -30)
	dateTo := now

	if c.QueryParam("date_from") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_from")); err == nil {
			dateFrom = parsed
		}
	}
	if c.QueryParam("date_to") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_to")); err == nil {
			dateTo = parsed
		}
	}

	stats, err := h.cognitiveGameService.GetGameStats(ctx, "", &dateFrom, &dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار بازی‌های شناختی", err)
	}

	scoreTrend, err := h.cognitiveGameService.GetScoreTrend(ctx, "", dateFrom, dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت روند نمرات", err)
	}

	timeAnalysis, err := h.cognitiveGameService.GetTimeAnalysis(ctx, "")
	if err != nil {
		return h.InternalError(c, "خطا در تحلیل زمان", err)
	}

	return response.OK(c, "الگوهای شناختی دریافت شد", map[string]interface{}{
		"date_range":    map[string]string{"from": dateFrom.Format("2006-01-02"), "to": dateTo.Format("2006-01-02")},
		"stats":         stats,
		"score_trend":   scoreTrend,
		"time_analysis": timeAnalysis,
	})
}

func (h *ReportHandler) GetMoodTrends(c echo.Context) error {
	ctx := c.Request().Context()
	now := time.Now()
	dateFrom := now.AddDate(0, 0, -30)
	dateTo := now

	if c.QueryParam("date_from") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_from")); err == nil {
			dateFrom = parsed
		}
	}
	if c.QueryParam("date_to") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_to")); err == nil {
			dateTo = parsed
		}
	}

	stats, err := h.moodTrackerService.GetMoodStats(ctx, "", &dateFrom, &dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار خلق", err)
	}

	trend, err := h.moodTrackerService.GetMoodTrend(ctx, "", dateFrom, dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت روند خلق", err)
	}

	activityEffectiveness, err := h.moodTrackerService.GetActivityEffectiveness(ctx, "")
	if err != nil {
		return h.InternalError(c, "خطا در دریافت اثربخشی فعالیت‌ها", err)
	}

	return response.OK(c, "روندهای خلق دریافت شد", map[string]interface{}{
		"date_range":             map[string]string{"from": dateFrom.Format("2006-01-02"), "to": dateTo.Format("2006-01-02")},
		"stats":                  stats,
		"mood_trend":             trend,
		"activity_effectiveness": activityEffectiveness,
	})
}

func (h *ReportHandler) GetEngagement(c echo.Context) error {
	ctx := c.Request().Context()
	now := time.Now()
	dateFrom := now.AddDate(0, 0, -30)
	dateTo := now

	if c.QueryParam("date_from") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_from")); err == nil {
			dateFrom = parsed
		}
	}
	if c.QueryParam("date_to") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_to")); err == nil {
			dateTo = parsed
		}
	}

	userEngagement, err := h.userService.GetUserEngagementStats(ctx, dateFrom, dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار مشارکت", err)
	}

	breathingStats, err := h.breathingService.GetSessionStats(ctx, "", &dateFrom, &dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار جلسات تنفس", err)
	}

	return response.OK(c, "آمار مشارکت دریافت شد", map[string]interface{}{
		"date_range":         map[string]string{"from": dateFrom.Format("2006-01-02"), "to": dateTo.Format("2006-01-02")},
		"user_engagement":    userEngagement,
		"breathing_activity": breathingStats,
	})
}

func (h *ReportHandler) GetWeeklyProgress(c echo.Context) error {
	ctx := c.Request().Context()
	now := time.Now()
	dateFrom := now.AddDate(0, 0, -30)
	dateTo := now

	if c.QueryParam("date_from") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_from")); err == nil {
			dateFrom = parsed
		}
	}
	if c.QueryParam("date_to") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_to")); err == nil {
			dateTo = parsed
		}
	}

	weeklyStats := []map[string]interface{}{}
	stressStats, err := h.stressEventService.GetStressStats(ctx, "", &dateFrom, &dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار استرس", err)
	}

	moodStats, err := h.moodTrackerService.GetMoodStats(ctx, "", &dateFrom, &dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار خلق", err)
	}

	cognitiveStats, err := h.cognitiveGameService.GetGameStats(ctx, "", &dateFrom, &dateTo)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار شناختی", err)
	}

	return response.OK(c, "پیشرفت هفتگی دریافت شد", map[string]interface{}{
		"date_range":     map[string]string{"from": dateFrom.Format("2006-01-02"), "to": dateTo.Format("2006-01-02")},
		"weekly_reports": weeklyStats,
		"stress":         stressStats,
		"mood":           moodStats,
		"cognitive":      cognitiveStats,
	})
}

func (h *ReportHandler) ExportData(c echo.Context) error {
	ctx := c.Request().Context()
	now := time.Now()
	dateFrom := now.AddDate(0, 0, -30)
	dateTo := now
	userID := c.QueryParam("user_id")

	if c.QueryParam("date_from") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_from")); err == nil {
			dateFrom = parsed
		}
	}
	if c.QueryParam("date_to") != "" {
		if parsed, err := time.Parse("2006-01-02", c.QueryParam("date_to")); err == nil {
			dateTo = parsed
		}
	}

	users, err := h.userService.ExportUsers(ctx, &dateFrom, &dateTo, userID)
	if err != nil {
		return h.InternalError(c, "خطا در خروجی داده‌ها", err)
	}

	stressStats, _ := h.stressEventService.GetStressStats(ctx, userID, &dateFrom, &dateTo)
	moodStats, _ := h.moodTrackerService.GetMoodStats(ctx, userID, &dateFrom, &dateTo)
	cognitiveStats, _ := h.cognitiveGameService.GetGameStats(ctx, userID, &dateFrom, &dateTo)

	return response.OK(c, "خروجی داده‌ها دریافت شد", map[string]interface{}{
		"date_range":       map[string]string{"from": dateFrom.Format("2006-01-02"), "to": dateTo.Format("2006-01-02")},
		"users":            users,
		"stress_analytics": stressStats,
		"mood_analytics":   moodStats,
		"cognitive_stats":  cognitiveStats,
	})
}

func (h *ReportHandler) GetWeeklyStats(c echo.Context) error {
	userID := h.GetUserID(c)
	if userID == "" {
		return h.Unauthorized(c, "شناسه کاربر یافت نشد")
	}

	stats, err := h.weeklyReportService.GetWeeklyStats(c.Request().Context(), userID)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت آمار هفتگی", err)
	}

	return response.OK(c, "آمار هفتگی با موفقیت دریافت شد", stats)
}

func (h *ReportHandler) BadRequest(c echo.Context, message string) error {
	return response.BadRequest(c, message, "")
}

func (h *ReportHandler) NotFound(c echo.Context, resource string) error {
	return response.Error(c, http.StatusNotFound, resource+" یافت نشد", "")
}

func (h *ReportHandler) InternalError(c echo.Context, message string, err error) error {
	return response.InternalServerError(c, message, err.Error())
}

func (h *ReportHandler) Unauthorized(c echo.Context, message string) error {
	return response.Unauthorized(c, message, "")
}

func (h *ReportHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}
