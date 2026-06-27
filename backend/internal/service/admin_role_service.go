package service

import (
	"context"
	"errors"
	"fmt"

	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type AdminRoleService struct {
	adminRoleRepo *repository.AdminRoleRepository
}

func NewAdminRoleService(adminRoleRepo *repository.AdminRoleRepository) *AdminRoleService {
	return &AdminRoleService{
		adminRoleRepo: adminRoleRepo,
	}
}

func (s *AdminRoleService) CreateAdminRole(ctx context.Context, req *schemas.AdminRoleCreateRequest) (*schemas.AdminRoleResponse, error) {
	if req.RoleName == "" {
		return nil, errors.New("نام نقش الزامی است")
	}

	existing, err := s.adminRoleRepo.GetByRoleName(ctx, req.RoleName)
	if err == nil && existing != nil {
		return nil, errors.New("این نام نقش قبلاً ثبت شده است")
	}

	adminRole := &models.AdminRole{
		RoleName:    req.RoleName,
		Description: req.Description,
		Permissions: req.Permissions,
	}

	if err := s.adminRoleRepo.Create(ctx, adminRole); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد نقش: %w", err)
	}

	return s.toAdminRoleResponse(adminRole), nil
}

func (s *AdminRoleService) GetAdminRoleById(ctx context.Context, roleId string) (*schemas.AdminRoleResponse, error) {
	role, err := s.adminRoleRepo.GetByID(ctx, roleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("نقش مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت نقش: %w", err)
	}

	return s.toAdminRoleResponse(role), nil
}

func (s *AdminRoleService) GetAllRoles(ctx context.Context, req *schemas.AdminRoleListRequest) (*schemas.AdminRoleListResponse, error) {
	filterFunc := s.buildAdminRoleFilters(req)

	roles, total, err := s.adminRoleRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت نقش‌ها: %w", err)
	}

	roleResponses := make([]schemas.AdminRoleResponse, len(roles))
	for i, role := range roles {
		roleResponses[i] = *s.toAdminRoleResponse(&role)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.AdminRoleListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.AdminRoleResponse]{
			Data:     roleResponses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *AdminRoleService) UpdateAdminRole(ctx context.Context, id string, req *schemas.AdminRoleUpdateRequest) (*schemas.AdminRoleResponse, error) {
	adminRole, err := s.adminRoleRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("نقش مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت نقش: %w", err)
	}

	if req.Description != nil {
		adminRole.Description = *req.Description
	}

	if req.Permissions != nil {
		adminRole.Permissions = *req.Permissions
	}

	if err := s.adminRoleRepo.Update(ctx, adminRole); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی نقش: %w", err)
	}

	return s.toAdminRoleResponse(adminRole), nil
}

func (s *AdminRoleService) DeleteAdminRole(ctx context.Context, id string) error {
	_, err := s.adminRoleRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("نقش مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت نقش: %w", err)
	}

	if err := s.adminRoleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف نقش: %w", err)
	}

	return nil
}

func (s *AdminRoleService) buildAdminRoleFilters(req *schemas.AdminRoleListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.Search != "" {
			db = db.Where("role_name LIKE ? OR description LIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *AdminRoleService) toAdminRoleResponse(adminRole *models.AdminRole) *schemas.AdminRoleResponse {
	return &schemas.AdminRoleResponse{
		ID:          adminRole.ID.String(),
		RoleName:    adminRole.RoleName,
		Description: adminRole.Description,
		Permissions: adminRole.Permissions,
		CreatedAt:   adminRole.CreatedAt,
		UpdatedAt:   adminRole.UpdatedAt,
	}
}
