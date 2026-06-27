package middleware

import (
	"net/http"
	"strings"

	"psychology-backend/internal/service"

	"github.com/labstack/echo/v4"
)

type JWTMiddleware struct {
	JWTService *service.JWTService
}

func NewJWTMiddleware(jwtService *service.JWTService) *JWTMiddleware {
	return &JWTMiddleware{JWTService: jwtService}
}

func (m *JWTMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Missing authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid authorization format. Use: Bearer <token>",
			})
		}

		claims, err := m.JWTService.ValidateToken(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid or expired token",
			})
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.UserRole)

		return next(c)
	}
}

func (m *JWTMiddleware) RequireUserRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("user_role").(string)
		if !ok || role != "user" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "User access required",
			})
		}
		return next(c)
	}
}

func (m *JWTMiddleware) RequireAdminRole(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("user_role").(string)
		if !ok || role != "admin" {
			return c.JSON(http.StatusForbidden, map[string]string{
				"error": "Admin access required",
			})
		}
		return next(c)
	}
}

func (m *JWTMiddleware) OptionalAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return next(c)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return next(c)
		}

		claims, err := m.JWTService.ValidateToken(tokenString)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("user_role", claims.UserRole)
		}

		return next(c)
	}
}
