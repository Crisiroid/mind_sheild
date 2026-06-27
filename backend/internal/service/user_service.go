package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserService handles business logic for User operations
// This is a COMPLETE EXAMPLE service that demonstrates:
// - CRUD operations (Create, Read, Update, Delete)
// - Pagination and filtering
// - Business logic and validation
// - Statistics and analytics
// - Error handling patterns
// - Repository integration
// - Schema transformation (model ↔ response)
type UserService struct {
	userRepo    *repository.UserRepository
	settingRepo *repository.UserSettingRepository
}

// NewUserService creates a new UserService instance
// Dependencies are injected here (repositories only - NO direct DB access)
func NewUserService(userRepo *repository.UserRepository, settingRepo *repository.UserSettingRepository) *UserService {
	return &UserService{
		userRepo:    userRepo,
		settingRepo: settingRepo,
	}
}

// ============================================================================
// CREATE OPERATIONS
// ============================================================================

// CreateUser handles the business logic for creating a new user
// This demonstrates:
// - Input validation
// - Business rule enforcement
// - Transaction handling (if needed)
// - Model creation from request schema
func (s *UserService) CreateUser(ctx context.Context, req *schemas.UserCreateRequest) (*schemas.UserResponse, error) {
	// Step 1: Validate business rules
	if err := s.validatePhoneNumber(req.PhoneNumber); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Step 2: Check if user already exists
	existingUser, err := s.userRepo.GetByPhoneNumber(ctx, req.PhoneNumber)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this phone number already exists")
	}

	// Step 3: Create the model from request
	now := time.Now()
	user := &models.User{
		PhoneNumber:      req.PhoneNumber,
		RegistrationDate: now,
		LastLogin:        &now,
		LoginCount:       1,
		AndroidVersion:   req.AndroidVersion,
		AppVersion:       req.AppVersion,
	}

	// Step 4: Save to database
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Step 5: Create default settings for the user
	defaultSettings := &models.UserSetting{
		UserID:               user.ID.String(),
		NotificationEnabled:  true,
		VibrationEnabled:     true,
		Language:             "fa",
		FontSize:             "medium",
		Theme:                "light",
		CrisisAlertThreshold: 3,
	}
	if err := s.settingRepo.Create(ctx, defaultSettings); err != nil {
		// Log the error but don't fail - user was created successfully
		// In production, you might want to use a proper logger here
		fmt.Printf("Warning: failed to create default settings for user %s: %v\n", user.ID, err)
	}

	// Step 6: Transform model to response schema
	response := s.toUserResponse(user)
	return response, nil
}

// ============================================================================
// READ OPERATIONS
// ============================================================================

// GetUserByID retrieves a user by their ID
// Demonstrates:
// - Simple retrieval
// - Error handling for not found cases
// - Model to response transformation
func (s *UserService) GetUserByID(ctx context.Context, id string) (*schemas.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return s.toUserResponse(user), nil
}

// GetUserByPhoneNumber retrieves a user by phone number
func (s *UserService) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (*schemas.UserResponse, error) {
	user, err := s.userRepo.GetByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return s.toUserResponse(user), nil
}

// ListUsers retrieves users with pagination and filters
// Demonstrates:
// - Complex filtering logic
// - Pagination handling
// - Building filter functions dynamically
// - Batch transformation of models to responses
func (s *UserService) ListUsers(ctx context.Context, req *schemas.UserListRequest) (*schemas.UserListResponse, error) {
	// Build the filter function based on request parameters
	filterFunc := s.buildUserFilters(req)

	// Get data from repository
	users, total, err := s.userRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Transform models to responses
	userResponses := make([]schemas.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.toUserResponse(&user)
	}

	// Calculate pagination metadata
	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	// Build and return the response
	return &schemas.UserListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.UserResponse]{
			Data:     userResponses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

// ============================================================================
// UPDATE OPERATIONS
// ============================================================================

// UpdateUser handles updating user information
// Demonstrates:
// - Partial updates (only update provided fields)
// - Fetching existing record
// - Validation before update
// - Selective field updates
func (s *UserService) UpdateUser(ctx context.Context, id string, req *schemas.UserUpdateRequest) (*schemas.UserResponse, error) {
	// Step 1: Get existing user
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Step 2: Update only the fields that were provided
	if req.CloudSyncEnabled != nil {
		user.CloudSyncEnabled = *req.CloudSyncEnabled
	}

	if req.DoNotDisturbEnabled != nil {
		user.DoNotDisturbEnabled = *req.DoNotDisturbEnabled
	}

	if req.DNDStartTime != nil {
		user.DNDStartTime = req.DNDStartTime
	}

	if req.DNDEndTime != nil {
		user.DNDEndTime = req.DNDEndTime
	}

	// Step 3: Validate business rules (if needed)
	if user.DNDStartTime != nil && user.DNDEndTime != nil {
		if user.DNDEndTime.Before(*user.DNDStartTime) {
			return nil, errors.New("DND end time must be after start time")
		}
	}

	// Step 4: Save updated user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.toUserResponse(user), nil
}

// AcceptAgreement handles user accepting the terms agreement
// Demonstrates:
// - Specific business operation
// - Timestamp tracking
// - Single field update
func (s *UserService) AcceptAgreement(ctx context.Context, userID string) (*schemas.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	now := time.Now()
	user.AgreementAccepted = true
	user.AgreementAcceptedAt = &now

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update agreement: %w", err)
	}

	return s.toUserResponse(user), nil
}

// UpdateLoginInfo updates user's login information
// Demonstrates:
// - Tracking login metrics
// - Incremental updates
func (s *UserService) UpdateLoginInfo(ctx context.Context, userID string, androidVersion, appVersion string) (*schemas.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	now := time.Now()
	user.LastLogin = &now
	user.LoginCount++

	if androidVersion != "" {
		user.AndroidVersion = androidVersion
	}
	if appVersion != "" {
		user.AppVersion = appVersion
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update login info: %w", err)
	}

	return s.toUserResponse(user), nil
}

// ============================================================================
// DELETE OPERATIONS
// ============================================================================

// DeleteUser handles user deletion
// Demonstrates:
// - Soft delete vs hard delete considerations
// - Cascading deletions (if needed)
func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	// Check if user exists
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Delete the user
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ============================================================================
// STATISTICS AND ANALYTICS
// ============================================================================

// GetUserStats returns comprehensive user statistics
// Demonstrates:
// - Aggregating data from repository
// - Complex business metrics
// - Response schema usage
func (s *UserService) GetUserStats(ctx context.Context) (*schemas.UserStatsResponse, error) {
	stats, err := s.userRepo.GetUserStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	return stats, nil
}

// GetUserActivityTrend returns user registration trends
func (s *UserService) GetUserActivityTrend(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	trends, err := s.userRepo.GetUserActivityTrend(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity trend: %w", err)
	}

	return trends, nil
}

// GetLoginAnalytics returns login frequency and patterns
func (s *UserService) GetLoginAnalytics(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	analytics, err := s.userRepo.GetLoginAnalytics(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get login analytics: %w", err)
	}

	return analytics, nil
}

// GetAgreementStats returns agreement acceptance statistics
func (s *UserService) GetAgreementStats(ctx context.Context) (int64, int64, float64, error) {
	totalUsers, agreedUsers, rate, err := s.userRepo.GetAgreementStats(ctx)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get agreement stats: %w", err)
	}

	return totalUsers, agreedUsers, rate, nil
}

// GetAppVersionDistribution returns distribution of users by app version
func (s *UserService) GetAppVersionDistribution(ctx context.Context) ([]schemas.DistributionStats, error) {
	distribution, err := s.userRepo.GetAppVersionDistribution(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get app version distribution: %w", err)
	}

	return distribution, nil
}

// GetInactiveUsers returns users who haven't logged in recently
func (s *UserService) GetInactiveUsers(ctx context.Context, daysThreshold int) ([]schemas.UserResponse, error) {
	users, err := s.userRepo.GetInactiveUsers(ctx, daysThreshold)
	if err != nil {
		return nil, fmt.Errorf("failed to get inactive users: %w", err)
	}

	responses := make([]schemas.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *s.toUserResponse(&user)
	}

	return responses, nil
}

// GetUserEngagementStats returns user engagement statistics
func (s *UserService) GetUserEngagementStats(ctx context.Context, dateFrom, dateTo time.Time) (*schemas.EngagementStatsResponse, error) {
	stats, err := s.userRepo.GetUserEngagementStats(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get engagement stats: %w", err)
	}

	return stats, nil
}

// ExportUsers exports user data based on filters
func (s *UserService) ExportUsers(ctx context.Context, dateFrom, dateTo *time.Time, userID string) ([]schemas.UserResponse, error) {
	users, err := s.userRepo.ExportUsers(ctx, dateFrom, dateTo, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to export users: %w", err)
	}

	responses := make([]schemas.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *s.toUserResponse(&user)
	}

	return responses, nil
}

// ============================================================================
// HELPER METHODS
// ============================================================================

// validatePhoneNumber validates phone number format
// Business rule: Iranian phone numbers must start with 09
func (s *UserService) validatePhoneNumber(phoneNumber string) error {
	if len(phoneNumber) != 11 {
		return errors.New("شماره تماس باید 11 رقم باشد")
	}

	if phoneNumber[0] != '0' || phoneNumber[1] != '9' {
		return errors.New("شماره تماس باید با 09 شروع شود")
	}

	// Check if all characters are digits
	for _, ch := range phoneNumber {
		if ch < '0' || ch > '9' {
			return errors.New("شماره تماس نباید حاوی حرف باشد")
		}
	}

	return nil
}

// buildUserFilters creates a filter function based on request parameters
// This pattern allows dynamic query building
func (s *UserService) buildUserFilters(req *schemas.UserListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Apply date range filter
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		// Apply search filter on phone number
		if req.Search != "" {
			db = db.Where("phone_number LIKE ?", "%"+req.Search+"%")
		}

		// Apply agreement filter
		if req.AgreementAccepted != nil {
			db = db.Where("agreement_accepted = ?", *req.AgreementAccepted)
		}

		// Apply sorting
		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			// Default sorting
			db = db.Order("created_at desc")
		}

		return db
	}
}

// toUserResponse transforms a User model to UserResponse schema
// This is a critical pattern: keep models and schemas separate
func (s *UserService) toUserResponse(user *models.User) *schemas.UserResponse {
	response := &schemas.UserResponse{
		ID:                  user.ID.String(),
		PhoneNumber:         user.PhoneNumber,
		RegistrationDate:    user.RegistrationDate,
		LastLogin:           user.LastLogin,
		LoginCount:          user.LoginCount,
		AgreementAccepted:   user.AgreementAccepted,
		AgreementAcceptedAt: user.AgreementAcceptedAt,
		CloudSyncEnabled:    user.CloudSyncEnabled,
		DoNotDisturbEnabled: user.DoNotDisturbEnabled,
		DNDStartTime:        user.DNDStartTime,
		DNDEndTime:          user.DNDEndTime,
		AndroidVersion:      user.AndroidVersion,
		AppVersion:          user.AppVersion,
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           user.UpdatedAt,
	}

	// Include settings if loaded
	if user.Settings.ID != uuid.Nil {
		response.Settings = &schemas.UserSettingResponse{
			ID:                   user.Settings.ID.String(),
			UserID:               user.Settings.UserID,
			NotificationEnabled:  user.Settings.NotificationEnabled,
			VibrationEnabled:     user.Settings.VibrationEnabled,
			Language:             user.Settings.Language,
			FontSize:             user.Settings.FontSize,
			Theme:                user.Settings.Theme,
			CrisisAlertThreshold: user.Settings.CrisisAlertThreshold,
			CreatedAt:            user.Settings.CreatedAt,
			UpdatedAt:            user.Settings.UpdatedAt,
		}
	}

	return response
}
