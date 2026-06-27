package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type AdminUserService struct {
	adminUserRepo *repository.AdminUserRepository
	adminRoleRepo *repository.AdminRoleRepository
}

func NewAdminUserService(adminUserRepo *repository.AdminUserRepository, adminRoleRepo *repository.AdminRoleRepository) *AdminUserService {
	return &AdminUserService{
		adminUserRepo: adminUserRepo,
		adminRoleRepo: adminRoleRepo,
	}
}

func (s *AdminUserService) CreateAdminUser(ctx context.Context, req *schemas.AdminCreateRequest) (*schemas.AdminUserResponse, error) {
	if req.Username == "" {
		return nil, errors.New("نام کاربری الزامی است")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("رمز عبور باید حداقل 6 کاراکتر باشد")
	}

	_, err := s.adminUserRepo.GetByUsername(ctx, req.Username)
	if err == nil {
		return nil, errors.New("این نام کاربری قبلاً ثبت شده است")
	}

	_, err = s.adminUserRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("این ایمیل قبلاً ثبت شده است")
	}

	if req.RoleID != "" {
		_, err := s.adminRoleRepo.GetByID(ctx, req.RoleID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("نقش مورد نظر یافت نشد")
			}
			return nil, fmt.Errorf("خطا در بررسی نقش: %w", err)
		}
	}

	adminUser := &models.AdminUser{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
		FullName:     req.FullName,
		RoleID:       req.RoleID,
		IsActive:     req.IsActive,
	}

	if err := s.adminUserRepo.Create(ctx, adminUser); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد مدیر: %w", err)
	}

	return s.toAdminUserResponse(adminUser), nil
}

func (s *AdminUserService) GetAdminUserById(ctx context.Context, adminUserId string) (*schemas.AdminUserResponse, error) {
	user, err := s.adminUserRepo.GetByID(ctx, adminUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("مدیر مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت مدیر: %w", err)
	}

	return s.toAdminUserResponse(user), nil
}

func (s *AdminUserService) GetAdminUserByUsername(ctx context.Context, adminUserName string) (*schemas.AdminUserResponse, error) {
	adminUser, err := s.adminUserRepo.GetByUsername(ctx, adminUserName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("مدیر مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت مدیر: %w", err)
	}
	return s.toAdminUserResponse(adminUser), nil
}

func (s *AdminUserService) GetAdminUserByEmail(ctx context.Context, adminUserEmail string) (*schemas.AdminUserResponse, error) {
	adminUser, err := s.adminUserRepo.GetByEmail(ctx, adminUserEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("مدیر مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت مدیر: %w", err)
	}
	return s.toAdminUserResponse(adminUser), nil
}

func (s *AdminUserService) GetAllAdmins(ctx context.Context, req *schemas.AdminListRequest) (*schemas.AdminListResponse, error) {
	filterFunc := s.buildAdminUserFilters(req)

	admins, total, err := s.adminUserRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت مدیران: %w", err)
	}

	adminUserResponses := make([]schemas.AdminUserResponse, len(admins))
	for i, adminUser := range admins {
		adminUserResponses[i] = *s.toAdminUserResponse(&adminUser)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.AdminListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.AdminUserResponse]{
			Data:     adminUserResponses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *AdminUserService) UpdateAdminUser(ctx context.Context, id string, req *schemas.AdminUpdateRequest) (*schemas.AdminUserResponse, error) {
	adminUser, err := s.adminUserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("مدیر مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت مدیر: %w", err)
	}

	if req.Email != nil {
		existing, err := s.adminUserRepo.GetByEmail(ctx, *req.Email)
		if err == nil && existing.ID.String() != id {
			return nil, errors.New("این ایمیل قبلاً توسط مدیر دیگری استفاده شده است")
		}
		adminUser.Email = *req.Email
	}

	if req.FullName != nil {
		adminUser.FullName = *req.FullName
	}

	if req.RoleID != nil {
		if *req.RoleID != "" {
			_, err := s.adminRoleRepo.GetByID(ctx, *req.RoleID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("نقش مورد نظر یافت نشد")
				}
				return nil, fmt.Errorf("خطا در بررسی نقش: %w", err)
			}
		}
		adminUser.RoleID = *req.RoleID
	}

	if req.IsActive != nil {
		adminUser.IsActive = *req.IsActive
	}

	if err := s.adminUserRepo.Update(ctx, adminUser); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی مدیر: %w", err)
	}

	return s.toAdminUserResponse(adminUser), nil
}

func (s *AdminUserService) UpdateAdminLogin(ctx context.Context, userID string) (*schemas.AdminUserResponse, error) {
	adminUser, err := s.adminUserRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت مدیر: %w", err)
	}

	adminUser.LastLogin = "now"

	if err := s.adminUserRepo.Update(ctx, adminUser); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی ورود: %w", err)
	}

	return s.toAdminUserResponse(adminUser), nil
}

func (s *AdminUserService) DeleteAdminUser(ctx context.Context, id string) error {
	_, err := s.adminUserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("مدیر مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت مدیر: %w", err)
	}

	if err := s.adminUserRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف مدیر: %w", err)
	}

	return nil
}

func (s *AdminUserService) DeactivateAdminUser(ctx context.Context, id string) (*schemas.AdminUserResponse, error) {
	adminUser, err := s.adminUserRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("مدیر مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت مدیر: %w", err)
	}

	adminUser.IsActive = false

	if err := s.adminUserRepo.Update(ctx, adminUser); err != nil {
		return nil, fmt.Errorf("خطا در غیرفعال کردن مدیر: %w", err)
	}

	return s.toAdminUserResponse(adminUser), nil
}

func (s *AdminUserService) GetAdminStats(ctx context.Context) (int64, int64, error) {
	totalAdmins, activeAdmins, err := s.adminUserRepo.GetAdminStats(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("خطا در دریافت آمار: %w", err)
	}

	return totalAdmins, activeAdmins, nil
}

func (s *AdminUserService) GetAdminActivityLog(ctx context.Context, dateFrom, dateTo time.Time) ([]repository.AdminActivityLog, error) {
	logs, err := s.adminUserRepo.GetAdminActivityLog(ctx, dateFrom, dateTo)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت لاگ فعالیت: %w", err)
	}

	return logs, nil
}

func (s *AdminUserService) GetAdminsByRole(ctx context.Context) ([]repository.RoleDistribution, error) {
	distributions, err := s.adminUserRepo.GetAdminsByRole(ctx)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت توزیع نقش‌ها: %w", err)
	}

	return distributions, nil
}

func (s *AdminUserService) GetInactiveAdmins(ctx context.Context, daysThreshold int) ([]schemas.AdminUserResponse, error) {
	admins, err := s.adminUserRepo.GetInactiveAdmins(ctx, daysThreshold)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت مدیران غیرفعال: %w", err)
	}

	responses := make([]schemas.AdminUserResponse, len(admins))
	for i, admin := range admins {
		responses[i] = *s.toAdminUserResponse(&admin)
	}

	return responses, nil
}

func (s *AdminUserService) GetAdminsByActiveStatus(ctx context.Context, isActive bool) ([]schemas.AdminUserResponse, error) {
	admins, err := s.adminUserRepo.GetAdminsByActiveStatus(ctx, isActive)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت مدیران: %w", err)
	}

	responses := make([]schemas.AdminUserResponse, len(admins))
	for i, admin := range admins {
		responses[i] = *s.toAdminUserResponse(&admin)
	}

	return responses, nil
}

func (s *AdminUserService) buildAdminUserFilters(req *schemas.AdminListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.Search != "" {
			db = db.Where("username LIKE ? OR email LIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
		}

		if req.RoleID != nil {
			db = db.Where("role_id = ?", *req.RoleID)
		}

		if req.IsActive != nil {
			db = db.Where("is_active = ?", *req.IsActive)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *AdminUserService) toAdminUserResponse(adminUser *models.AdminUser) *schemas.AdminUserResponse {
	return &schemas.AdminUserResponse{
		ID:        adminUser.ID.String(),
		Username:  adminUser.Username,
		Email:     adminUser.Email,
		FullName:  adminUser.FullName,
		RoleID:    adminUser.RoleID,
		IsActive:  adminUser.IsActive,
		LastLogin: adminUser.LastLogin,
		CreatedAt: adminUser.CreatedAt,
		UpdatedAt: adminUser.UpdatedAt,
	}
}
