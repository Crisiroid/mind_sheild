package service

import (
	"context"
	"errors"
	"fmt"

	"psychology-backend/internal/interfaces"
	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type RoleValueService struct {
	roleValueRepo interfaces.RoleValueRepositoryInterface
}

func NewRoleValueService(roleValueRepo interfaces.RoleValueRepositoryInterface) *RoleValueService {
	return &RoleValueService{
		roleValueRepo: roleValueRepo,
	}
}

func (s *RoleValueService) CreateRoleValue(ctx context.Context, req *schemas.RoleValueCreateRequest) (*schemas.RoleValueResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	if req.EntryText == "" {
		return nil, errors.New("متن ورودی الزامی است")
	}

	roleValue := &models.RoleAndValue{
		UserID:    req.UserID,
		EntryType: req.EntryType,
		EntryText: req.EntryText,
		DayNumber: req.DayNumber,
	}

	if err := s.roleValueRepo.Create(ctx, roleValue); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد نقش و ارزش: %w", err)
	}

	return s.toRoleValueResponse(roleValue), nil
}

func (s *RoleValueService) GetRoleValueById(ctx context.Context, id string) (*schemas.RoleValueResponse, error) {
	roleValue, err := s.roleValueRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("نقش و ارزش مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت نقش و ارزش: %w", err)
	}

	return s.toRoleValueResponse(roleValue), nil
}

func (s *RoleValueService) GetAllRoleValues(ctx context.Context, req *schemas.RoleValueListRequest) (*schemas.RoleValueListResponse, error) {
	filterFunc := s.buildRoleValueFilters(req)

	roleValues, total, err := s.roleValueRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت نقش‌ها و ارزش‌ها: %w", err)
	}

	responses := make([]schemas.RoleValueResponse, len(roleValues))
	for i, roleValue := range roleValues {
		responses[i] = *s.toRoleValueResponse(&roleValue)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.RoleValueListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.RoleValueResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *RoleValueService) UpdateRoleValue(ctx context.Context, id string, req *schemas.RoleValueUpdateRequest) (*schemas.RoleValueResponse, error) {
	roleValue, err := s.roleValueRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("نقش و ارزش مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت نقش و ارزش: %w", err)
	}

	if req.EntryText != nil {
		roleValue.EntryText = *req.EntryText
	}

	if err := s.roleValueRepo.Update(ctx, roleValue); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی نقش و ارزش: %w", err)
	}

	return s.toRoleValueResponse(roleValue), nil
}

func (s *RoleValueService) DeleteRoleValue(ctx context.Context, id string) error {
	_, err := s.roleValueRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("نقش و ارزش مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت نقش و ارزش: %w", err)
	}

	if err := s.roleValueRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف نقش و ارزش: %w", err)
	}

	return nil
}

func (s *RoleValueService) buildRoleValueFilters(req *schemas.RoleValueListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.EntryType != nil {
			db = db.Where("entry_type = ?", *req.EntryType)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *RoleValueService) toRoleValueResponse(roleValue *models.RoleAndValue) *schemas.RoleValueResponse {
	return &schemas.RoleValueResponse{
		ID:          roleValue.ID.String(),
		UserID:      roleValue.UserID,
		EntryType:   roleValue.EntryType,
		EntryText:   roleValue.EntryText,
		CreatedDate: roleValue.CreatedDate,
		DayNumber:   roleValue.DayNumber,
		CreatedAt:   roleValue.CreatedAt,
		UpdatedAt:   roleValue.UpdatedAt,
	}
}
