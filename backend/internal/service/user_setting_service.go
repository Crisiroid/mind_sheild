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

type UserSettingService struct {
	settingRepo *repository.UserSettingRepository
}

func NewUserSettingService(settingRepo *repository.UserSettingRepository) *UserSettingService {
	return &UserSettingService{
		settingRepo: settingRepo,
	}
}

func (s *UserSettingService) GetUserSettingById(ctx context.Context, id string) (*schemas.UserSettingResponse, error) {
	setting, err := s.settingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تنظیمات کاربر مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت تنظیمات کاربر: %w", err)
	}

	return s.toUserSettingResponse(setting), nil
}

func (s *UserSettingService) GetUserSettingByUserId(ctx context.Context, userID string) (*schemas.UserSettingResponse, error) {
	setting, err := s.settingRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تنظیمات کاربر مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت تنظیمات کاربر: %w", err)
	}

	return s.toUserSettingResponse(setting), nil
}

func (s *UserSettingService) GetAllUserSettings(ctx context.Context, req *schemas.UserSettingListRequest) (*schemas.UserSettingListResponse, error) {
	filterFunc := s.buildUserSettingFilters(req)

	settings, total, err := s.settingRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت تنظیمات کاربران: %w", err)
	}

	responses := make([]schemas.UserSettingResponse, len(settings))
	for i, setting := range settings {
		responses[i] = *s.toUserSettingResponse(&setting)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.UserSettingListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.UserSettingResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *UserSettingService) UpdateUserSetting(ctx context.Context, id string, req *schemas.UserSettingUpdateRequest) (*schemas.UserSettingResponse, error) {
	setting, err := s.settingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تنظیمات کاربر مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت تنظیمات کاربر: %w", err)
	}

	if req.NotificationEnabled != nil {
		setting.NotificationEnabled = *req.NotificationEnabled
	}

	if req.VibrationEnabled != nil {
		setting.VibrationEnabled = *req.VibrationEnabled
	}

	if req.Language != nil {
		setting.Language = *req.Language
	}

	if req.FontSize != nil {
		setting.FontSize = *req.FontSize
	}

	if req.Theme != nil {
		setting.Theme = *req.Theme
	}

	if req.CrisisAlertThreshold != nil {
		setting.CrisisAlertThreshold = *req.CrisisAlertThreshold
	}

	if err := s.settingRepo.Update(ctx, setting); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی تنظیمات کاربر: %w", err)
	}

	return s.toUserSettingResponse(setting), nil
}

func (s *UserSettingService) UpdateUserSettingByUserId(ctx context.Context, userID string, req *schemas.UserSettingUpdateRequest) (*schemas.UserSettingResponse, error) {
	setting, err := s.settingRepo.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تنظیمات کاربر مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت تنظیمات کاربر: %w", err)
	}

	if req.NotificationEnabled != nil {
		setting.NotificationEnabled = *req.NotificationEnabled
	}

	if req.VibrationEnabled != nil {
		setting.VibrationEnabled = *req.VibrationEnabled
	}

	if req.Language != nil {
		setting.Language = *req.Language
	}

	if req.FontSize != nil {
		setting.FontSize = *req.FontSize
	}

	if req.Theme != nil {
		setting.Theme = *req.Theme
	}

	if req.CrisisAlertThreshold != nil {
		setting.CrisisAlertThreshold = *req.CrisisAlertThreshold
	}

	if err := s.settingRepo.Update(ctx, setting); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی تنظیمات کاربر: %w", err)
	}

	return s.toUserSettingResponse(setting), nil
}

func (s *UserSettingService) DeleteUserSetting(ctx context.Context, id string) error {
	_, err := s.settingRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("تنظیمات کاربر مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت تنظیمات کاربر: %w", err)
	}

	if err := s.settingRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف تنظیمات کاربر: %w", err)
	}

	return nil
}

func (s *UserSettingService) buildUserSettingFilters(req *schemas.UserSettingListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.Language != nil {
			db = db.Where("language = ?", *req.Language)
		}

		if req.Theme != nil {
			db = db.Where("theme = ?", *req.Theme)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *UserSettingService) toUserSettingResponse(setting *models.UserSetting) *schemas.UserSettingResponse {
	return &schemas.UserSettingResponse{
		ID:                   setting.ID.String(),
		UserID:               setting.UserID,
		NotificationEnabled:  setting.NotificationEnabled,
		VibrationEnabled:     setting.VibrationEnabled,
		Language:             setting.Language,
		FontSize:             setting.FontSize,
		Theme:                setting.Theme,
		CrisisAlertThreshold: setting.CrisisAlertThreshold,
		CreatedAt:            setting.CreatedAt,
		UpdatedAt:            setting.UpdatedAt,
	}
}
