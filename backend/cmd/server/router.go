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

	setupAuthRoutes(e, authHandler, adminHandler)

	setupUserRoutes(e, authHandler, userHandler, calendarHandler, emotionHandler, breathingHandler, cognitiveHandler, mentalMustHandler, negativeThoughtHandler, mindCourtHandler, conflictExerciseHandler, moodTrackerHandler, roleValueHandler, skyThoughtHandler, mindfulnessHandler, reportHandler, mediaContentHandler, jwtMiddleware)

	setupAdminRoutes(e, authHandler, adminHandler, userHandler, calendarHandler, emotionHandler, breathingHandler, cognitiveHandler, mentalMustHandler, negativeThoughtHandler, mindCourtHandler, conflictExerciseHandler, moodTrackerHandler, roleValueHandler, skyThoughtHandler, mindfulnessHandler, reportHandler, mediaContentHandler, jwtMiddleware)
}

func setupAuthRoutes(e *echo.Echo, authHandler *handler.AuthHandler, adminHandler *handler.AdminHandler) {
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

	adminApi.GET("/admin/me", adminHandler.GetAdminProfile)
	adminApi.PUT("/admin/me", adminHandler.UpdateAdminProfile)

	setupAdminUserRoutes(adminApi, userHandler)

	setupAdminManagementRoutes(adminApi, adminHandler)

	setupAdminRoleRoutes(adminApi, adminHandler)

	adminApi.GET("/admin/logs", adminHandler.ListSystemLogs)
	adminApi.GET("/admin/logs/:id", adminHandler.GetSystemLogByID)

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
	g.GET("/users", userHandler.ListUsers)
	g.DELETE("/users/:id", userHandler.DeleteUser)
	g.GET("/users/by-phone", userHandler.GetUserByPhoneNumber)
	g.GET("/users/stats", userHandler.GetUserStats)
	g.GET("/users/activity-trend", userHandler.GetUserActivityTrend)
	g.GET("/users/login-analytics", userHandler.GetLoginAnalytics)
	g.GET("/users/agreement-stats", userHandler.GetAgreementStats)
	g.GET("/users/app-version-distribution", userHandler.GetAppVersionDistribution)
	g.GET("/users/inactive", userHandler.GetInactiveUsers)
	g.GET("/users/engagement", userHandler.GetUserEngagement)
	g.GET("/users/export", userHandler.ExportUsers)
}

func setupAdminManagementRoutes(g *echo.Group, adminHandler *handler.AdminHandler) {
	g.POST("/admin/users", adminHandler.CreateAdminUser)
	g.GET("/admin/users/:id", adminHandler.GetAdminUserByID)
	g.GET("/admin/users", adminHandler.ListAdminUsers)
	g.PUT("/admin/users/:id", adminHandler.UpdateAdminUser)
	g.DELETE("/admin/users/:id", adminHandler.DeleteAdminUser)
	g.POST("/admin/users/:id/deactivate", adminHandler.DeactivateAdminUser)
}

func setupAdminRoleRoutes(g *echo.Group, adminHandler *handler.AdminHandler) {
	g.POST("/admin/roles", adminHandler.CreateAdminRole)
	g.GET("/admin/roles/:id", adminHandler.GetAdminRoleByID)
	g.GET("/admin/roles", adminHandler.ListAdminRoles)
	g.PUT("/admin/roles/:id", adminHandler.UpdateAdminRole)
	g.DELETE("/admin/roles/:id", adminHandler.DeleteAdminRole)
}

func setupAdminReportRoutes(g *echo.Group, reportHandler *handler.ReportHandler) {
	g.POST("/admin/reports", reportHandler.CreateUserReport)
	g.GET("/admin/reports/:id", reportHandler.GetUserReportByID)
	g.GET("/admin/reports", reportHandler.ListUserReports)
	g.DELETE("/admin/reports/:id", reportHandler.DeleteUserReport)

	g.GET("/reports/dashboard", reportHandler.GetDashboard)
	g.GET("/reports/user-activity", reportHandler.GetUserActivity)
	g.GET("/reports/stress-analytics", reportHandler.GetStressAnalytics)
	g.GET("/reports/body-tension", reportHandler.GetBodyTensionReport)
	g.GET("/reports/cognitive-patterns", reportHandler.GetCognitivePatterns)
	g.GET("/reports/mood-trends", reportHandler.GetMoodTrends)
	g.GET("/reports/engagement", reportHandler.GetEngagement)
	g.GET("/reports/weekly-progress", reportHandler.GetWeeklyProgress)
	g.GET("/reports/export", reportHandler.ExportData)
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
	g.GET("/calendars", calendarHandler.ListCalendarEntries)
	g.GET("/emotion-interactions", emotionHandler.ListEmotionInteractions)
	g.GET("/stress-events", emotionHandler.ListStressEvents)
	g.GET("/body-tension-maps", emotionHandler.ListBodyTensionMaps)
	g.GET("/breathing-sessions", breathingHandler.ListBreathingSessions)
	g.GET("/cognitive-games", cognitiveHandler.ListCognitiveGames)
	g.GET("/mental-musts", mentalMustHandler.ListMentalMusts)
	g.GET("/negative-thoughts", negativeThoughtHandler.ListNegativeThoughts)
	g.GET("/mind-court-evidence", mindCourtHandler.ListMindCourtEvidence)
	g.GET("/conflict-exercises", conflictExerciseHandler.ListConflictExercises)
	g.GET("/mood-tracker", moodTrackerHandler.ListMoodTrackers)
	g.GET("/roles-values", roleValueHandler.ListRolesValues)
	g.GET("/sky-thoughts", skyThoughtHandler.ListSkyThoughts)
	g.GET("/mindful-timers", mindfulnessHandler.ListMindfulTimers)
	g.GET("/acceptance-exercises", mindfulnessHandler.ListAcceptanceExercises)
	g.GET("/reports/weekly", reportHandler.ListWeeklyReports)
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
	g.GET("/media/weekly", mediaContentHandler.ListMediaContent)
	g.GET("/media/weekly/:id", mediaContentHandler.GetMediaContentByID)
	g.GET("/media/weekly/by-week/:week_number", mediaContentHandler.GetMediaContentByWeek)
	g.GET("/media/weekly/:id/download", mediaContentHandler.DownloadMediaContent)
}

func setupAdminMediaContentRoutes(g *echo.Group, mediaContentHandler *handler.WeeklyMediaContentHandler) {
	g.POST("/media/weekly", mediaContentHandler.UploadMediaContent)
	g.GET("/media/weekly", mediaContentHandler.ListMediaContent)
	g.GET("/media/weekly/:id", mediaContentHandler.GetMediaContentByID)
	g.PUT("/media/weekly/:id", mediaContentHandler.UpdateMediaContent)
	g.DELETE("/media/weekly/:id", mediaContentHandler.DeleteMediaContent)
	g.GET("/media/weekly/by-week/:week_number", mediaContentHandler.GetMediaContentByWeek)
	g.GET("/media/weekly/:id/download", mediaContentHandler.DownloadMediaContent)
}
