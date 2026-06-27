package main

import (
	"net/http"
	"time"

	"psychology-backend/internal/handler"
	"psychology-backend/internal/middleware"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	adminHandler *handler.AdminHandler,
	calendarHandler *handler.CalendarHandler,
	emotionHandler *handler.EmotionHandler,
	breathingHandler *handler.BreathingHandler,
	cognitiveHandler *handler.CognitiveHandler,
	mentalMustHandler *handler.MentalMustHandler,
	negativeThoughtHandler *handler.NegativeThoughtHandler,
	mindCourtHandler *handler.MindCourtHandler,
	conflictExerciseHandler *handler.ConflictExerciseHandler,
	moodTrackerHandler *handler.MoodTrackerHandler,
	roleValueHandler *handler.RoleValueHandler,
	skyThoughtHandler *handler.SkyThoughtHandler,
	mindfulnessHandler *handler.MindfulnessHandler,
	reportHandler *handler.ReportHandler,
	mediaContentHandler *handler.WeeklyMediaContentHandler,
	jwtMiddleware *middleware.JWTMiddleware,
) {
	e.GET("/health", healthCheckHandler)
	e.GET("/api/v1/public/health", publicHealthCheckHandler)

	setupAuthRoutes(e, authHandler)

	setupUserRoutes(e, authHandler, userHandler, calendarHandler, emotionHandler, breathingHandler, cognitiveHandler, mentalMustHandler, negativeThoughtHandler, mindCourtHandler, conflictExerciseHandler, moodTrackerHandler, roleValueHandler, skyThoughtHandler, mindfulnessHandler, reportHandler, mediaContentHandler, jwtMiddleware)

	setupAdminRoutes(e, authHandler, adminHandler, userHandler, calendarHandler, emotionHandler, breathingHandler, cognitiveHandler, mentalMustHandler, negativeThoughtHandler, mindCourtHandler, conflictExerciseHandler, moodTrackerHandler, roleValueHandler, skyThoughtHandler, mindfulnessHandler, reportHandler, mediaContentHandler, jwtMiddleware)
}

func setupAuthRoutes(e *echo.Echo, authHandler *handler.AuthHandler) {
	e.POST(schemas.RouteUserRegister, authHandler.UserRegister)
	e.POST(schemas.RouteUserLogin, authHandler.UserLogin)
	e.POST(schemas.RouteUserRefreshToken, authHandler.UserRefreshToken)
	e.POST(schemas.RouteAdminLogin, authHandler.AdminLogin)
	e.POST(schemas.RouteAdminRefreshToken, authHandler.AdminRefreshToken)
}

func setupUserRoutes(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	calendarHandler *handler.CalendarHandler,
	emotionHandler *handler.EmotionHandler,
	breathingHandler *handler.BreathingHandler,
	cognitiveHandler *handler.CognitiveHandler,
	mentalMustHandler *handler.MentalMustHandler,
	negativeThoughtHandler *handler.NegativeThoughtHandler,
	mindCourtHandler *handler.MindCourtHandler,
	conflictExerciseHandler *handler.ConflictExerciseHandler,
	moodTrackerHandler *handler.MoodTrackerHandler,
	roleValueHandler *handler.RoleValueHandler,
	skyThoughtHandler *handler.SkyThoughtHandler,
	mindfulnessHandler *handler.MindfulnessHandler,
	reportHandler *handler.ReportHandler,
	mediaContentHandler *handler.WeeklyMediaContentHandler,
	jwtMiddleware *middleware.JWTMiddleware,
) {
	userAuth := e.Group("")
	userAuth.Use(jwtMiddleware.Authenticate, jwtMiddleware.RequireUserRole)
	userAuth.POST(schemas.RouteUserLogout, authHandler.UserLogout)
	userAuth.POST(schemas.RouteUserChangePassword, authHandler.UserChangePassword)

	userApi := e.Group("/api/v1")
	userApi.Use(jwtMiddleware.Authenticate, jwtMiddleware.RequireUserRole)

	userApi.GET(schemas.RouteUserProfile, userHandler.GetUserProfile)
	userApi.PUT(schemas.RouteUserProfile, userHandler.UpdateUserProfile)
	userApi.POST("/api/v1/users/me/sync", userHandler.SyncUserData)
	userApi.POST("/api/v1/users/me/agreement", userHandler.AcceptAgreement)

	userApi.POST(schemas.RouteUsers, userHandler.CreateUser)
	userApi.GET(schemas.RouteUserByID, userHandler.GetUserByID)
	userApi.PUT(schemas.RouteUserByID, userHandler.UpdateUser)
	userApi.POST("/api/v1/users/:id/accept-agreement", userHandler.AcceptAgreement)
	userApi.POST("/api/v1/users/update-login-info", userHandler.UpdateLoginInfo)

	setupCalendarRoutes(userApi, calendarHandler)

	setupEmotionRoutes(userApi, emotionHandler)

	setupExerciseRoutes(userApi, breathingHandler, cognitiveHandler, mentalMustHandler, negativeThoughtHandler, mindCourtHandler, conflictExerciseHandler, moodTrackerHandler, roleValueHandler, skyThoughtHandler)

	setupMindfulnessRoutes(userApi, mindfulnessHandler)

	setupUserReportRoutes(userApi, reportHandler)

	setupMediaContentRoutes(userApi, mediaContentHandler)
}

func setupAdminRoutes(
	e *echo.Echo,
	authHandler *handler.AuthHandler,
	adminHandler *handler.AdminHandler,
	userHandler *handler.UserHandler,
	calendarHandler *handler.CalendarHandler,
	emotionHandler *handler.EmotionHandler,
	breathingHandler *handler.BreathingHandler,
	cognitiveHandler *handler.CognitiveHandler,
	mentalMustHandler *handler.MentalMustHandler,
	negativeThoughtHandler *handler.NegativeThoughtHandler,
	mindCourtHandler *handler.MindCourtHandler,
	conflictExerciseHandler *handler.ConflictExerciseHandler,
	moodTrackerHandler *handler.MoodTrackerHandler,
	roleValueHandler *handler.RoleValueHandler,
	skyThoughtHandler *handler.SkyThoughtHandler,
	mindfulnessHandler *handler.MindfulnessHandler,
	reportHandler *handler.ReportHandler,
	mediaContentHandler *handler.WeeklyMediaContentHandler,
	jwtMiddleware *middleware.JWTMiddleware,
) {
	adminAuth := e.Group("")
	adminAuth.Use(jwtMiddleware.Authenticate, jwtMiddleware.RequireAdminRole)
	adminAuth.POST(schemas.RouteAdminLogout, authHandler.AdminLogout)
	adminAuth.POST(schemas.RouteAdminChangePassword, authHandler.AdminChangePassword)

	adminApi := e.Group("/api/v1")
	adminApi.Use(jwtMiddleware.Authenticate, jwtMiddleware.RequireAdminRole)

	adminApi.GET(schemas.RouteAdminProfile, adminHandler.GetAdminProfile)
	adminApi.PUT(schemas.RouteAdminProfile, adminHandler.UpdateAdminProfile)

	setupAdminUserRoutes(adminApi, userHandler)

	setupAdminManagementRoutes(adminApi, adminHandler)

	setupAdminRoleRoutes(adminApi, adminHandler)

	adminApi.GET(schemas.RouteAdminLogs, adminHandler.ListSystemLogs)
	adminApi.GET(schemas.RouteAdminLogsByID, adminHandler.GetSystemLogByID)

	setupAdminReportRoutes(adminApi, reportHandler)

	setupAdminDataRoutes(adminApi, calendarHandler, emotionHandler, breathingHandler, cognitiveHandler, mentalMustHandler, negativeThoughtHandler, mindCourtHandler, conflictExerciseHandler, moodTrackerHandler, roleValueHandler, skyThoughtHandler, mindfulnessHandler, reportHandler)

	setupAdminMediaContentRoutes(adminApi, mediaContentHandler)
}

func setupCalendarRoutes(g *echo.Group, calendarHandler *handler.CalendarHandler) {
	g.POST(schemas.RouteCalendars, calendarHandler.CreateCalendarEntry)
	g.GET(schemas.RouteCalendarByID, calendarHandler.GetCalendarEntryByID)
	g.GET(schemas.RouteCalendars, calendarHandler.ListCalendarEntries)
	g.PUT(schemas.RouteCalendarByID, calendarHandler.UpdateCalendarEntry)
	g.DELETE(schemas.RouteCalendarByID, calendarHandler.DeleteCalendarEntry)
	g.GET("/api/v1/calendars/stats/completion", calendarHandler.GetCompletionStats)
	g.GET("/api/v1/calendars/stats/progress", calendarHandler.GetDayRangeProgress)
	g.GET("/api/v1/calendars/stats/streak", calendarHandler.GetStreakAnalysis)
}

func setupEmotionRoutes(g *echo.Group, emotionHandler *handler.EmotionHandler) {
	g.POST(schemas.RouteEmotionInteractions, emotionHandler.CreateEmotionInteraction)
	g.GET(schemas.RouteEmotionInteractionByID, emotionHandler.GetEmotionInteractionByID)
	g.GET(schemas.RouteEmotionInteractions, emotionHandler.ListEmotionInteractions)
	g.PUT(schemas.RouteEmotionInteractionByID, emotionHandler.UpdateEmotionInteraction)
	g.DELETE(schemas.RouteEmotionInteractionByID, emotionHandler.DeleteEmotionInteraction)

	g.POST(schemas.RouteStressEvents, emotionHandler.CreateStressEvent)
	g.GET(schemas.RouteStressEventByID, emotionHandler.GetStressEventByID)
	g.GET(schemas.RouteStressEvents, emotionHandler.ListStressEvents)
	g.PUT(schemas.RouteStressEventByID, emotionHandler.UpdateStressEvent)
	g.DELETE(schemas.RouteStressEventByID, emotionHandler.DeleteStressEvent)

	g.POST(schemas.RouteBodyTensionMaps, emotionHandler.CreateBodyTensionMap)
	g.GET(schemas.RouteBodyTensionMapByID, emotionHandler.GetBodyTensionMapByID)
	g.GET(schemas.RouteBodyTensionMaps, emotionHandler.ListBodyTensionMaps)
	g.PUT(schemas.RouteBodyTensionMapByID, emotionHandler.UpdateBodyTensionMap)
	g.DELETE(schemas.RouteBodyTensionMapByID, emotionHandler.DeleteBodyTensionMap)
}

func setupExerciseRoutes(
	g *echo.Group,
	breathingHandler *handler.BreathingHandler,
	cognitiveHandler *handler.CognitiveHandler,
	mentalMustHandler *handler.MentalMustHandler,
	negativeThoughtHandler *handler.NegativeThoughtHandler,
	mindCourtHandler *handler.MindCourtHandler,
	conflictExerciseHandler *handler.ConflictExerciseHandler,
	moodTrackerHandler *handler.MoodTrackerHandler,
	roleValueHandler *handler.RoleValueHandler,
	skyThoughtHandler *handler.SkyThoughtHandler,
) {
	g.POST(schemas.RouteBreathingSessions, breathingHandler.CreateBreathingSession)
	g.GET(schemas.RouteBreathingSessionByID, breathingHandler.GetBreathingSessionByID)
	g.GET(schemas.RouteBreathingSessions, breathingHandler.ListBreathingSessions)
	g.PUT(schemas.RouteBreathingSessionByID, breathingHandler.UpdateBreathingSession)
	g.DELETE(schemas.RouteBreathingSessionByID, breathingHandler.DeleteBreathingSession)

	g.POST(schemas.RouteCognitiveGames, cognitiveHandler.CreateCognitiveGame)
	g.GET(schemas.RouteCognitiveGameByID, cognitiveHandler.GetCognitiveGameByID)
	g.GET(schemas.RouteCognitiveGames, cognitiveHandler.ListCognitiveGames)
	g.PUT(schemas.RouteCognitiveGameByID, cognitiveHandler.UpdateCognitiveGame)
	g.DELETE(schemas.RouteCognitiveGameByID, cognitiveHandler.DeleteCognitiveGame)

	g.POST(schemas.RouteMentalMusts, mentalMustHandler.CreateMentalMust)
	g.GET(schemas.RouteMentalMustByID, mentalMustHandler.GetMentalMustByID)
	g.GET(schemas.RouteMentalMusts, mentalMustHandler.ListMentalMusts)
	g.PUT(schemas.RouteMentalMustByID, mentalMustHandler.UpdateMentalMust)
	g.DELETE(schemas.RouteMentalMustByID, mentalMustHandler.DeleteMentalMust)

	g.POST(schemas.RouteNegativeThoughts, negativeThoughtHandler.CreateNegativeThought)
	g.GET(schemas.RouteNegativeThoughtByID, negativeThoughtHandler.GetNegativeThoughtByID)
	g.GET(schemas.RouteNegativeThoughts, negativeThoughtHandler.ListNegativeThoughts)
	g.PUT(schemas.RouteNegativeThoughtByID, negativeThoughtHandler.UpdateNegativeThought)
	g.DELETE(schemas.RouteNegativeThoughtByID, negativeThoughtHandler.DeleteNegativeThought)

	g.POST(schemas.RouteMindCourtEvidence, mindCourtHandler.CreateMindCourtEvidence)
	g.GET(schemas.RouteMindCourtEvidenceByID, mindCourtHandler.GetMindCourtEvidenceByID)
	g.GET(schemas.RouteMindCourtEvidence, mindCourtHandler.ListMindCourtEvidence)
	g.PUT(schemas.RouteMindCourtEvidenceByID, mindCourtHandler.UpdateMindCourtEvidence)
	g.DELETE(schemas.RouteMindCourtEvidenceByID, mindCourtHandler.DeleteMindCourtEvidence)

	g.POST(schemas.RouteConflictExercises, conflictExerciseHandler.CreateConflictExercise)
	g.GET(schemas.RouteConflictExerciseByID, conflictExerciseHandler.GetConflictExerciseByID)
	g.GET(schemas.RouteConflictExercises, conflictExerciseHandler.ListConflictExercises)
	g.PUT(schemas.RouteConflictExerciseByID, conflictExerciseHandler.UpdateConflictExercise)
	g.DELETE(schemas.RouteConflictExerciseByID, conflictExerciseHandler.DeleteConflictExercise)

	g.POST(schemas.RouteMoodTracker, moodTrackerHandler.CreateMoodTracker)
	g.GET(schemas.RouteMoodTrackerByID, moodTrackerHandler.GetMoodTrackerByID)
	g.GET(schemas.RouteMoodTracker, moodTrackerHandler.ListMoodTrackers)
	g.PUT(schemas.RouteMoodTrackerByID, moodTrackerHandler.UpdateMoodTracker)
	g.DELETE(schemas.RouteMoodTrackerByID, moodTrackerHandler.DeleteMoodTracker)

	g.POST(schemas.RouteRolesValues, roleValueHandler.CreateRoleValue)
	g.GET(schemas.RouteRolesValuesByID, roleValueHandler.GetRoleValueByID)
	g.GET(schemas.RouteRolesValues, roleValueHandler.ListRolesValues)
	g.PUT(schemas.RouteRolesValuesByID, roleValueHandler.UpdateRoleValue)
	g.DELETE(schemas.RouteRolesValuesByID, roleValueHandler.DeleteRoleValue)

	g.POST(schemas.RouteSkyThoughts, skyThoughtHandler.CreateSkyThought)
	g.GET(schemas.RouteSkyThoughtByID, skyThoughtHandler.GetSkyThoughtByID)
	g.GET(schemas.RouteSkyThoughts, skyThoughtHandler.ListSkyThoughts)
	g.PUT(schemas.RouteSkyThoughtByID, skyThoughtHandler.UpdateSkyThought)
	g.DELETE(schemas.RouteSkyThoughtByID, skyThoughtHandler.DeleteSkyThought)
}

func setupMindfulnessRoutes(g *echo.Group, mindfulnessHandler *handler.MindfulnessHandler) {
	g.POST(schemas.RouteMindfulTimers, mindfulnessHandler.CreateMindfulTimer)
	g.GET(schemas.RouteMindfulTimerByID, mindfulnessHandler.GetMindfulTimerByID)
	g.GET(schemas.RouteMindfulTimers, mindfulnessHandler.ListMindfulTimers)
	g.PUT(schemas.RouteMindfulTimerByID, mindfulnessHandler.UpdateMindfulTimer)
	g.DELETE(schemas.RouteMindfulTimerByID, mindfulnessHandler.DeleteMindfulTimer)

	g.POST(schemas.RouteAcceptanceExercises, mindfulnessHandler.CreateAcceptanceExercise)
	g.GET(schemas.RouteAcceptanceExerciseByID, mindfulnessHandler.GetAcceptanceExerciseByID)
	g.GET(schemas.RouteAcceptanceExercises, mindfulnessHandler.ListAcceptanceExercises)
	g.PUT(schemas.RouteAcceptanceExerciseByID, mindfulnessHandler.UpdateAcceptanceExercise)
	g.DELETE(schemas.RouteAcceptanceExerciseByID, mindfulnessHandler.DeleteAcceptanceExercise)
}

func setupUserReportRoutes(g *echo.Group, reportHandler *handler.ReportHandler) {
	g.POST("/api/v1/reports/weekly", reportHandler.CreateWeeklyReport)
	g.GET("/api/v1/reports/weekly/:id", reportHandler.GetWeeklyReportByID)
	g.GET("/api/v1/reports/weekly", reportHandler.ListWeeklyReports)
	g.PUT("/api/v1/reports/weekly/:id", reportHandler.UpdateWeeklyReport)
	g.DELETE("/api/v1/reports/weekly/:id", reportHandler.DeleteWeeklyReport)
	g.GET(schemas.RouteReportsDashboard, reportHandler.GetDashboard)
	g.GET(schemas.RouteReportsUserActivity, reportHandler.GetUserActivity)
	g.GET(schemas.RouteReportsStressAnalytics, reportHandler.GetStressAnalytics)
	g.GET(schemas.RouteReportsBodyTension, reportHandler.GetBodyTensionReport)
	g.GET(schemas.RouteReportsCognitivePatterns, reportHandler.GetCognitivePatterns)
	g.GET(schemas.RouteReportsMoodTrends, reportHandler.GetMoodTrends)
	g.GET(schemas.RouteReportsEngagement, reportHandler.GetEngagement)
	g.GET(schemas.RouteReportsWeeklyProgress, reportHandler.GetWeeklyProgress)
	g.GET(schemas.RouteReportsExport, reportHandler.ExportData)
	g.GET("/api/v1/reports/weekly-stats", reportHandler.GetWeeklyStats)
}

func setupAdminUserRoutes(g *echo.Group, userHandler *handler.UserHandler) {
	g.GET(schemas.RouteUsers, userHandler.ListUsers)
	g.DELETE(schemas.RouteUserByID, userHandler.DeleteUser)
	g.GET("/api/v1/users/by-phone", userHandler.GetUserByPhoneNumber)
	g.GET("/api/v1/users/stats", userHandler.GetUserStats)
	g.GET("/api/v1/users/activity-trend", userHandler.GetUserActivityTrend)
	g.GET("/api/v1/users/login-analytics", userHandler.GetLoginAnalytics)
	g.GET("/api/v1/users/agreement-stats", userHandler.GetAgreementStats)
	g.GET("/api/v1/users/app-version-distribution", userHandler.GetAppVersionDistribution)
	g.GET("/api/v1/users/inactive", userHandler.GetInactiveUsers)
	g.GET("/api/v1/users/engagement", userHandler.GetUserEngagement)
	g.GET("/api/v1/users/export", userHandler.ExportUsers)
}

func setupAdminManagementRoutes(g *echo.Group, adminHandler *handler.AdminHandler) {
	g.POST(schemas.RouteAdminUsers, adminHandler.CreateAdminUser)
	g.GET(schemas.RouteAdminUserByID, adminHandler.GetAdminUserByID)
	g.GET(schemas.RouteAdminUsers, adminHandler.ListAdminUsers)
	g.PUT(schemas.RouteAdminUserByID, adminHandler.UpdateAdminUser)
	g.DELETE(schemas.RouteAdminUserByID, adminHandler.DeleteAdminUser)
	g.POST("/api/v1/admin/users/:id/deactivate", adminHandler.DeactivateAdminUser)
}

func setupAdminRoleRoutes(g *echo.Group, adminHandler *handler.AdminHandler) {
	g.POST(schemas.RouteAdminRoles, adminHandler.CreateAdminRole)
	g.GET(schemas.RouteAdminRoleByID, adminHandler.GetAdminRoleByID)
	g.GET(schemas.RouteAdminRoles, adminHandler.ListAdminRoles)
	g.PUT(schemas.RouteAdminRoleByID, adminHandler.UpdateAdminRole)
	g.DELETE(schemas.RouteAdminRoleByID, adminHandler.DeleteAdminRole)
}

func setupAdminReportRoutes(g *echo.Group, reportHandler *handler.ReportHandler) {
	g.POST(schemas.RouteAdminReports, reportHandler.CreateUserReport)
	g.GET(schemas.RouteAdminReportsByID, reportHandler.GetUserReportByID)
	g.GET(schemas.RouteAdminReports, reportHandler.ListUserReports)
	g.DELETE(schemas.RouteAdminReportsByID, reportHandler.DeleteUserReport)

	g.GET(schemas.RouteReportsDashboard, reportHandler.GetDashboard)
	g.GET(schemas.RouteReportsUserActivity, reportHandler.GetUserActivity)
	g.GET(schemas.RouteReportsStressAnalytics, reportHandler.GetStressAnalytics)
	g.GET(schemas.RouteReportsBodyTension, reportHandler.GetBodyTensionReport)
	g.GET(schemas.RouteReportsCognitivePatterns, reportHandler.GetCognitivePatterns)
	g.GET(schemas.RouteReportsMoodTrends, reportHandler.GetMoodTrends)
	g.GET(schemas.RouteReportsEngagement, reportHandler.GetEngagement)
	g.GET(schemas.RouteReportsWeeklyProgress, reportHandler.GetWeeklyProgress)
	g.GET(schemas.RouteReportsExport, reportHandler.ExportData)
}

func setupAdminDataRoutes(
	g *echo.Group,
	calendarHandler *handler.CalendarHandler,
	emotionHandler *handler.EmotionHandler,
	breathingHandler *handler.BreathingHandler,
	cognitiveHandler *handler.CognitiveHandler,
	mentalMustHandler *handler.MentalMustHandler,
	negativeThoughtHandler *handler.NegativeThoughtHandler,
	mindCourtHandler *handler.MindCourtHandler,
	conflictExerciseHandler *handler.ConflictExerciseHandler,
	moodTrackerHandler *handler.MoodTrackerHandler,
	roleValueHandler *handler.RoleValueHandler,
	skyThoughtHandler *handler.SkyThoughtHandler,
	mindfulnessHandler *handler.MindfulnessHandler,
	reportHandler *handler.ReportHandler,
) {
	g.GET(schemas.RouteCalendars, calendarHandler.ListCalendarEntries)
	g.GET(schemas.RouteEmotionInteractions, emotionHandler.ListEmotionInteractions)
	g.GET(schemas.RouteStressEvents, emotionHandler.ListStressEvents)
	g.GET(schemas.RouteBodyTensionMaps, emotionHandler.ListBodyTensionMaps)
	g.GET(schemas.RouteBreathingSessions, breathingHandler.ListBreathingSessions)
	g.GET(schemas.RouteCognitiveGames, cognitiveHandler.ListCognitiveGames)
	g.GET(schemas.RouteMentalMusts, mentalMustHandler.ListMentalMusts)
	g.GET(schemas.RouteNegativeThoughts, negativeThoughtHandler.ListNegativeThoughts)
	g.GET(schemas.RouteMindCourtEvidence, mindCourtHandler.ListMindCourtEvidence)
	g.GET(schemas.RouteConflictExercises, conflictExerciseHandler.ListConflictExercises)
	g.GET(schemas.RouteMoodTracker, moodTrackerHandler.ListMoodTrackers)
	g.GET(schemas.RouteRolesValues, roleValueHandler.ListRolesValues)
	g.GET(schemas.RouteSkyThoughts, skyThoughtHandler.ListSkyThoughts)
	g.GET(schemas.RouteMindfulTimers, mindfulnessHandler.ListMindfulTimers)
	g.GET(schemas.RouteAcceptanceExercises, mindfulnessHandler.ListAcceptanceExercises)
	g.GET("/api/v1/reports/weekly", reportHandler.ListWeeklyReports)
}

func healthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func publicHealthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "psychology-backend",
	})
}

func setupMediaContentRoutes(g *echo.Group, mediaContentHandler *handler.WeeklyMediaContentHandler) {
	g.GET(schemas.RouteWeeklyMediaContents, mediaContentHandler.ListMediaContent)
	g.GET(schemas.RouteWeeklyMediaContentByID, mediaContentHandler.GetMediaContentByID)
	g.GET(schemas.RouteWeeklyMediaByWeek, mediaContentHandler.GetMediaContentByWeek)
	g.GET(schemas.RouteMediaDownload, mediaContentHandler.DownloadMediaContent)
}

func setupAdminMediaContentRoutes(g *echo.Group, mediaContentHandler *handler.WeeklyMediaContentHandler) {
	g.POST(schemas.RouteWeeklyMediaContents, mediaContentHandler.UploadMediaContent)
	g.GET(schemas.RouteWeeklyMediaContents, mediaContentHandler.ListMediaContent)
	g.GET(schemas.RouteWeeklyMediaContentByID, mediaContentHandler.GetMediaContentByID)
	g.PUT(schemas.RouteWeeklyMediaContentByID, mediaContentHandler.UpdateMediaContent)
	g.DELETE(schemas.RouteWeeklyMediaContentByID, mediaContentHandler.DeleteMediaContent)
	g.GET(schemas.RouteWeeklyMediaByWeek, mediaContentHandler.GetMediaContentByWeek)
	g.GET(schemas.RouteMediaDownload, mediaContentHandler.DownloadMediaContent)
}
