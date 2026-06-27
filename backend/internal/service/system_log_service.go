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

type SystemLogService struct {
	systemLogRepo interfaces.SystemLogRepositoryInterface
}

func NewSystemLogService(systemLogRepo interfaces.SystemLogRepositoryInterface) *SystemLogService {
	return &SystemLogService{
		systemLogRepo: systemLogRepo,
	}
}

func (s *SystemLogService) GetSystemLogById(ctx context.Context, id string) (*schemas.SystemLogResponse, error) {
	log, err := s.systemLogRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("لاگ سیستم مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت لاگ سیستم: %w", err)
	}

	return s.toSystemLogResponse(log), nil
}

func (s *SystemLogService) GetAllSystemLogs(ctx context.Context, req *schemas.SystemLogListRequest) (*schemas.SystemLogListResponse, error) {
	filterFunc := s.buildSystemLogFilters(req)

	logs, total, err := s.systemLogRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت لاگ‌های سیستم: %w", err)
	}

	responses := make([]schemas.SystemLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = *s.toSystemLogResponse(&log)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.SystemLogListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.SystemLogResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *SystemLogService) DeleteSystemLog(ctx context.Context, id string) error {
	_, err := s.systemLogRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("لاگ سیستم مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت لاگ سیستم: %w", err)
	}

	if err := s.systemLogRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف لاگ سیستم: %w", err)
	}

	return nil
}

func (s *SystemLogService) buildSystemLogFilters(req *schemas.SystemLogListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.LogType != nil {
			db = db.Where("log_type = ?", *req.LogType)
		}

		if req.Severity != nil {
			db = db.Where("severity = ?", *req.Severity)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *SystemLogService) toSystemLogResponse(log *models.SystemLog) *schemas.SystemLogResponse {
	return &schemas.SystemLogResponse{
		ID:         log.ID.String(),
		LogType:    log.LogType,
		LogMessage: log.LogMessage,
		UserID:     log.UserID,
		Severity:   log.Severity,
		CreatedAt:  log.CreatedAt,
	}
}
