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
	"psychology-backend/internal/middleware"
	"psychology-backend/pkg/validator"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	log.Printf("Starting application in %s mode", cfg.App.Env)

	// Connect to database
	db, err := database.Connect(&cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run auto-migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Echo
	e := echo.New()

	// Set custom validator
	e.Validator = validator.NewCustomValidator()

	// Middleware
	e.Pre(echoMiddleware.RemoveTrailingSlash())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.CORS())
	e.Use(middleware.LoggerMiddleware)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// API routes
	api := e.Group("/api/v1")

	// Public routes
	public := api.Group("/public")
	public.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "psychology-backend",
		})
	})

	// Protected routes (will add JWT middleware later)
	// protected := api.Group("")
	// protected.Use(middleware.NewJWTMiddleware(cfg.JWT.Secret).Authenticate)

	// Start server
	go func() {
		addr := ":" + cfg.App.Port
		log.Printf("Server starting on %s", addr)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Shutting down the server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Close database connection
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
