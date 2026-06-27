package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"psychology-backend/internal/interfaces"
	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo interfaces.UserRepositoryInterface
}

func NewUserService(userRepo interfaces.UserRepositoryInterface) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *schemas.UserCreateRequest) (*schemas.UserResponse, error) {
	if err := s.validatePhoneNumber(req.PhoneNumber); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	existingUser, err := s.userRepo.GetByPhoneNumber(ctx, req.PhoneNumber)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this phone number already exists")
	}

	now := time.Now()
	user := &models.User{
		PhoneNumber:      req.PhoneNumber,
		RegistrationDate: now,
		LastLogin:        &now,
		LoginCount:       1,
		AndroidVersion:   req.AndroidVersion,
		AppVersion:       req.AppVersion,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	response := s.toUserResponse(user)
	return response, nil
}

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

func (s *UserService) ListUsers(ctx context.Context, req *schemas.UserListRequest) (*schemas.UserListResponse, error) {
	filterFunc := s.buildUserFilters(req)

	users, total, err := s.userRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	userResponses := make([]schemas.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = *s.toUserResponse(&user)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

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

func (s *UserService) UpdateUser(ctx context.Context, id string, req *schemas.UserUpdateRequest) (*schemas.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

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

	if user.DNDStartTime != nil && user.DNDEndTime != nil {
		if user.DNDEndTime.Before(*user.DNDStartTime) {
			return nil, errors.New("DND end time must be after start time")
		}
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.toUserResponse(user), nil
}

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

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (s *UserService) GetUserStats(ctx context.Context) (*schemas.UserStatsResponse, error) {
	stats, err := s.userRepo.GetUserStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	return stats, nil
}

func (s *UserService) GetUserActivityTrend(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	trends, err := s.userRepo.GetUserActivityTrend(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity trend: %w", err)
	}

	return trends, nil
}

func (s *UserService) GetLoginAnalytics(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	analytics, err := s.userRepo.GetLoginAnalytics(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get login analytics: %w", err)
	}

	return analytics, nil
}

func (s *UserService) GetAgreementStats(ctx context.Context) (int64, int64, float64, error) {
	totalUsers, agreedUsers, rate, err := s.userRepo.GetAgreementStats(ctx)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get agreement stats: %w", err)
	}

	return totalUsers, agreedUsers, rate, nil
}

func (s *UserService) GetAppVersionDistribution(ctx context.Context) ([]schemas.DistributionStats, error) {
	distribution, err := s.userRepo.GetAppVersionDistribution(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get app version distribution: %w", err)
	}

	return distribution, nil
}

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

func (s *UserService) GetUserEngagementStats(ctx context.Context, dateFrom, dateTo time.Time) (*schemas.EngagementStatsResponse, error) {
	stats, err := s.userRepo.GetUserEngagementStats(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("failed to get engagement stats: %w", err)
	}

	return stats, nil
}

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

func (s *UserService) validatePhoneNumber(phoneNumber string) error {
	if len(phoneNumber) != 11 {
		return errors.New("شماره تماس باید 11 رقم باشد")
	}

	if phoneNumber[0] != '0' || phoneNumber[1] != '9' {
		return errors.New("شماره تماس باید با 09 شروع شود")
	}

	for _, ch := range phoneNumber {
		if ch < '0' || ch > '9' {
			return errors.New("شماره تماس نباید حاوی حرف باشد")
		}
	}

	return nil
}

func (s *UserService) buildUserFilters(req *schemas.UserListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.Search != "" {
			db = db.Where("phone_number LIKE ?", "%"+req.Search+"%")
		}

		if req.AgreementAccepted != nil {
			db = db.Where("agreement_accepted = ?", *req.AgreementAccepted)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *UserService) toUserResponse(user *models.User) *schemas.UserResponse {
	return &schemas.UserResponse{
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
}

func (s *UserService) GetUserProfile(ctx context.Context, userID string) (*schemas.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("کاربر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت کاربر: %w", err)
	}

	return s.toUserResponse(user), nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, userID string, req *schemas.UserUpdateProfileRequest) (*schemas.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("کاربر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت کاربر: %w", err)
	}

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

	if user.DNDStartTime != nil && user.DNDEndTime != nil {
		if user.DNDEndTime.Before(*user.DNDStartTime) {
			return nil, errors.New("زمان پایان مزاحم نشوید باید بعد از زمان شروع باشد")
		}
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی کاربر: %w", err)
	}

	return s.toUserResponse(user), nil
}

func (s *UserService) SyncUserData(ctx context.Context, userID string, syncData map[string]interface{}) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("کاربر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت کاربر: %w", err)
	}

	user.CloudSyncEnabled = true
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("خطا در همگام‌سازی داده‌ها: %w", err)
	}

	return nil
}
