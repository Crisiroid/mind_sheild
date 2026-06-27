package service

import (
	"fmt"
	"time"

	"psychology-backend/pkg/schemas"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	SecretKey         string
	AccessExpiration  time.Duration
	RefreshExpiration time.Duration
}

func NewJWTService(secret string, accessExp, refreshExp time.Duration) *JWTService {
	return &JWTService{
		SecretKey:         secret,
		AccessExpiration:  accessExp,
		RefreshExpiration: refreshExp,
	}
}

func (s *JWTService) GenerateAccessToken(userID, userRole string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":   userID,
		"user_role": userRole,
		"exp":       now.Add(s.AccessExpiration).Unix(),
		"iat":       now.Unix(),
		"nbf":       now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

func (s *JWTService) GenerateRefreshToken() (string, time.Time, error) {
	token := uuid.New().String()
	expiry := time.Now().Add(s.RefreshExpiration)
	return token, expiry, nil
}

func (s *JWTService) ValidateToken(tokenString string) (*schemas.TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.SecretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_id in token")
	}

	userRole, ok := claims["user_role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid user_role in token")
	}

	var email string
	if e, ok := claims["email"].(string); ok {
		email = e
	}

	tokenClaims := &schemas.TokenClaims{
		UserID:   userID,
		UserRole: userRole,
		Email:    email,
		IssuedAt: int64(claims["iat"].(float64)),
	}

	return tokenClaims, nil
}

func (s *JWTService) GetAccessExpiration() time.Duration {
	return s.AccessExpiration
}
