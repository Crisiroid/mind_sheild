package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"psychology-backend/internal/config"
	"psychology-backend/internal/database"
	"psychology-backend/internal/handler"
	"psychology-backend/internal/interfaces"
	"psychology-backend/internal/middleware"
	"psychology-backend/internal/repository"
	"psychology-backend/internal/service"
	"psychology-backend/pkg/validator"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.LoadConfig()
	log.Printf("Starting application in %s mode", cfg.App.Env)

	db, err := database.Connect(&cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	e := echo.New()

	e.Validator = validator.NewCustomValidator()

	e.Pre(echoMiddleware.RemoveTrailingSlash())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.CORS())
	e.Use(middleware.LoggerMiddleware)

	var userRepo interfaces.UserRepositoryInterface = repository.NewUserRepository(db)
	var adminRepo interfaces.AdminUserRepositoryInterface = repository.NewAdminUserRepository(db)
	var adminRoleRepo interfaces.AdminRoleRepositoryInterface = repository.NewAdminRoleRepository(db)
	var dailyCalendarRepo interfaces.DailyCalendarRepositoryInterface = repository.NewDailyCalendarRepository(db)
	var breathingRepo interfaces.BreathingRepositoryInterface = repository.NewBreathingRepository(db)
	var emotionRepo interfaces.EmotionTriangleRepositoryInterface = repository.NewEmotionTriangleRepository(db)
	var stressRepo interfaces.StressEventRepositoryInterface = repository.NewStressEventRepository(db)
	var bodyTensionRepo interfaces.BodyTensionRepositoryInterface = repository.NewBodyTensionRepository(db)
	var cognitiveGameRepo interfaces.CognitiveGameRepositoryInterface = repository.NewCognitiveGameRepository(db)
	var mentalMustRepo interfaces.MentalMustRepositoryInterface = repository.NewMentalMustRepository(db)
	var negativeThoughtRepo interfaces.NegativeThoughtRepositoryInterface = repository.NewNegativeThoughtRepository(db)
	var mindCourtRepo interfaces.MindCourtRepositoryInterface = repository.NewMindCourtRepository(db)
	var conflictExerciseRepo interfaces.ConflictExerciseRepositoryInterface = repository.NewConflictExerciseRepository(db)
	var moodTrackerRepo interfaces.MoodTrackerRepositoryInterface = repository.NewMoodTrackerRepository(db)
	var roleValueRepo interfaces.RoleValueRepositoryInterface = repository.NewRoleValueRepository(db)
	var skyThoughtRepo interfaces.SkyThoughtRepositoryInterface = repository.NewSkyThoughtRepository(db)
	var mindfulTimerRepo interfaces.MindfulTimerRepositoryInterface = repository.NewMindfulTimerRepository(db)
	var acceptanceRepo interfaces.AcceptanceRepositoryInterface = repository.NewAcceptanceRepository(db)
	var userReportRepo interfaces.UserReportRepositoryInterface = repository.NewUserReportRepository(db)
	var weeklyReportRepo interfaces.WeeklyReportRepositoryInterface = repository.NewWeeklyReportRepository(db)
	var systemLogRepo interfaces.SystemLogRepositoryInterface = repository.NewSystemLogRepository(db)
	var weeklyMediaContentRepo interfaces.WeeklyMediaContentRepositoryInterface = repository.NewWeeklyMediaContentRepository(db)

	jwtService := service.NewJWTService(
		cfg.JWT.Secret,
		24*time.Hour,
		720*time.Hour,
	)
	passwordService := service.NewPasswordService(0)

	var authService interfaces.AuthServiceInterface = service.NewAuthService(userRepo, adminRepo, jwtService, passwordService)
	var userService interfaces.UserServiceInterface = service.NewUserService(userRepo)
	var adminUserService interfaces.AdminUserServiceInterface = service.NewAdminUserService(adminRepo, adminRoleRepo, passwordService)
	var adminRoleService interfaces.AdminRoleServiceInterface = service.NewAdminRoleService(adminRoleRepo)
	var dailyCalendarService interfaces.DailyCalendarServiceInterface = service.NewDailyCalendarService(dailyCalendarRepo)
	var breathingService interfaces.BreathingServiceInterface = service.NewBreathingService(breathingRepo)
	var emotionService interfaces.EmotionTriangleServiceInterface = service.NewEmotionTriangleService(emotionRepo)
	var stressService interfaces.StressEventServiceInterface = service.NewStressEventService(stressRepo)
	var bodyTensionService interfaces.BodyTensionServiceInterface = service.NewBodyTensionService(bodyTensionRepo)
	var cognitiveGameService interfaces.CognitiveGameServiceInterface = service.NewCognitiveGameService(cognitiveGameRepo)
	var mentalMustService interfaces.MentalMustServiceInterface = service.NewMentalMustService(mentalMustRepo)
	var negativeThoughtService interfaces.NegativeThoughtServiceInterface = service.NewNegativeThoughtService(negativeThoughtRepo)
	var mindCourtService interfaces.MindCourtServiceInterface = service.NewMindCourtService(mindCourtRepo)
	var conflictExerciseService interfaces.ConflictExerciseServiceInterface = service.NewConflictExerciseService(conflictExerciseRepo)
	var moodTrackerService interfaces.MoodTrackerServiceInterface = service.NewMoodTrackerService(moodTrackerRepo)
	var roleValueService interfaces.RoleValueServiceInterface = service.NewRoleValueService(roleValueRepo)
	var skyThoughtService interfaces.SkyThoughtServiceInterface = service.NewSkyThoughtService(skyThoughtRepo)
	var mindfulTimerService interfaces.MindfulTimerServiceInterface = service.NewMindfulTimerService(mindfulTimerRepo)
	var acceptanceService interfaces.AcceptanceServiceInterface = service.NewAcceptanceService(acceptanceRepo)
	var userReportService interfaces.UserReportServiceInterface = service.NewUserReportService(userReportRepo)
	var weeklyReportService interfaces.WeeklyReportServiceInterface = service.NewWeeklyReportService(weeklyReportRepo)
	var systemLogService interfaces.SystemLogServiceInterface = service.NewSystemLogService(systemLogRepo)
	var weeklyMediaContentService interfaces.WeeklyMediaContentServiceInterface = service.NewWeeklyMediaContentService(weeklyMediaContentRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	adminHandler := handler.NewAdminHandler(adminUserService, adminRoleService, systemLogService)
	calendarHandler := handler.NewCalendarHandler(dailyCalendarService)
	emotionHandler := handler.NewEmotionHandler(emotionService, stressService, bodyTensionService)

	breathingHandler := handler.NewBreathingHandler(breathingService)
	cognitiveHandler := handler.NewCognitiveHandler(cognitiveGameService)
	mentalMustHandler := handler.NewMentalMustHandler(mentalMustService)
	negativeThoughtHandler := handler.NewNegativeThoughtHandler(negativeThoughtService)
	mindCourtHandler := handler.NewMindCourtHandler(mindCourtService)
	conflictExerciseHandler := handler.NewConflictExerciseHandler(conflictExerciseService)
	moodTrackerHandler := handler.NewMoodTrackerHandler(moodTrackerService)
	roleValueHandler := handler.NewRoleValueHandler(roleValueService)
	skyThoughtHandler := handler.NewSkyThoughtHandler(skyThoughtService)

	mindfulnessHandler := handler.NewMindfulnessHandler(mindfulTimerService, acceptanceService)
	reportHandler := handler.NewReportHandler(userReportService, weeklyReportService, userService, stressService, bodyTensionService, cognitiveGameService, moodTrackerService, breathingService)

	uploadDirectory := config.GetEnv("UPLOAD_DIR", "./uploads")
	baseURL := config.GetEnv("BASE_URL", "http://localhost:8080")
	mediaContentHandler := handler.NewWeeklyMediaContentHandler(weeklyMediaContentService, uploadDirectory, baseURL)

	jwtMiddleware := middleware.NewJWTMiddleware(jwtService)

	SetupRoutes(e, authHandler, userHandler, adminHandler, calendarHandler, emotionHandler, breathingHandler, cognitiveHandler, mentalMustHandler, negativeThoughtHandler, mindCourtHandler, conflictExerciseHandler, moodTrackerHandler, roleValueHandler, skyThoughtHandler, mindfulnessHandler, reportHandler, mediaContentHandler, jwtMiddleware)

	go func() {
		addr := ":" + cfg.App.Port
		log.Printf("Server starting on %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Shutting down the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	if err := database.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
