package schemas

const (
	RouteUserRegister       = "/api/v1/auth/user/register"
	RouteUserLogin          = "/api/v1/auth/user/login"
	RouteUserRefreshToken   = "/api/v1/auth/user/refresh"
	RouteUserLogout         = "/api/v1/auth/user/logout"
	RouteUserChangePassword = "/api/v1/auth/user/change-password"

	RouteAdminLogin          = "/api/v1/auth/admin/login"
	RouteAdminRefreshToken   = "/api/v1/auth/admin/refresh"
	RouteAdminLogout         = "/api/v1/auth/admin/logout"
	RouteAdminChangePassword = "/api/v1/auth/admin/change-password"
)

const (
	RouteUsers       = "/api/v1/users"
	RouteUserByID    = "/api/v1/users/:id"
	RouteUserProfile = "/api/v1/users/me"
)

const (
	RouteCalendars    = "/api/v1/calendars"
	RouteCalendarByID = "/api/v1/calendars/:id"
)

const (
	RouteEmotionInteractions    = "/api/v1/emotion-interactions"
	RouteEmotionInteractionByID = "/api/v1/emotion-interactions/:id"

	RouteStressEvents    = "/api/v1/stress-events"
	RouteStressEventByID = "/api/v1/stress-events/:id"

	RouteBodyTensionMaps    = "/api/v1/body-tension-maps"
	RouteBodyTensionMapByID = "/api/v1/body-tension-maps/:id"
)

const (
	RouteBreathingSessions    = "/api/v1/breathing-sessions"
	RouteBreathingSessionByID = "/api/v1/breathing-sessions/:id"

	RouteCognitiveGames    = "/api/v1/cognitive-games"
	RouteCognitiveGameByID = "/api/v1/cognitive-games/:id"

	RouteMentalMusts    = "/api/v1/mental-musts"
	RouteMentalMustByID = "/api/v1/mental-musts/:id"

	RouteNegativeThoughts    = "/api/v1/negative-thoughts"
	RouteNegativeThoughtByID = "/api/v1/negative-thoughts/:id"

	RouteMindCourtEvidence     = "/api/v1/mind-court-evidence"
	RouteMindCourtEvidenceByID = "/api/v1/mind-court-evidence/:id"

	RouteConflictExercises    = "/api/v1/conflict-exercises"
	RouteConflictExerciseByID = "/api/v1/conflict-exercises/:id"

	RouteMoodTracker     = "/api/v1/mood-tracker"
	RouteMoodTrackerByID = "/api/v1/mood-tracker/:id"

	RouteRolesValues     = "/api/v1/roles-values"
	RouteRolesValuesByID = "/api/v1/roles-values/:id"

	RouteSkyThoughts    = "/api/v1/sky-thoughts"
	RouteSkyThoughtByID = "/api/v1/sky-thoughts/:id"
)

const (
	RouteMindfulTimers    = "/api/v1/mindful-timers"
	RouteMindfulTimerByID = "/api/v1/mindful-timers/:id"

	RouteAcceptanceExercises    = "/api/v1/acceptance-exercises"
	RouteAcceptanceExerciseByID = "/api/v1/acceptance-exercises/:id"
)

const (
	RouteReportsDashboard         = "/api/v1/reports/dashboard"
	RouteReportsUserActivity      = "/api/v1/reports/user-activity"
	RouteReportsStressAnalytics   = "/api/v1/reports/stress-analytics"
	RouteReportsBodyTension       = "/api/v1/reports/body-tension"
	RouteReportsCognitivePatterns = "/api/v1/reports/cognitive-patterns"
	RouteReportsMoodTrends        = "/api/v1/reports/mood-trends"
	RouteReportsEngagement        = "/api/v1/reports/engagement"
	RouteReportsWeeklyProgress    = "/api/v1/reports/weekly-progress"
	RouteReportsExport            = "/api/v1/reports/export"
)

const (
	RouteAdminUsers    = "/api/v1/admin/users"
	RouteAdminUserByID = "/api/v1/admin/users/:id"
	RouteAdminProfile  = "/api/v1/admin/me"
)

const (
	RouteAdminRoles    = "/api/v1/admin/roles"
	RouteAdminRoleByID = "/api/v1/admin/roles/:id"
)

const (
	RouteAdminReports     = "/api/v1/admin/reports"
	RouteAdminReportsByID = "/api/v1/admin/reports/:id"
	RouteAdminLogs        = "/api/v1/admin/logs"
	RouteAdminLogsByID    = "/api/v1/admin/logs/:id"
)

const (
	MethodGET    = "GET"
	MethodPOST   = "POST"
	MethodPUT    = "PUT"
	MethodPATCH  = "PATCH"
	MethodDELETE = "DELETE"
)
