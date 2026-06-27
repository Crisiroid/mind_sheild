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

type UserReportService struct {
	userReportRepo interfaces.UserReportRepositoryInterface
}

func NewUserReportService(userReportRepo interfaces.UserReportRepositoryInterface) *UserReportService {
	return &UserReportService{
		userReportRepo: userReportRepo,
	}
}

func (s *UserReportService) CreateUserReport(ctx context.Context, req *schemas.UserReportCreateRequest) (*schemas.UserReportResponse, error) {
	if req.ReportType == "" {
		return nil, errors.New("نوع گزارش الزامی است")
	}

	if req.ReportDate == "" {
		return nil, errors.New("تاریخ گزارش الزامی است")
	}

	report := &models.UserReport{
		ReportType:         req.ReportType,
		ReportDate:         req.ReportDate,
		TotalUsers:         req.TotalUsers,
		ActiveUsers:        req.ActiveUsers,
		AvgEngagementScore: req.AvgEngagementScore,
		CrisisAlertsCount:  req.CrisisAlertsCount,
		AnonymizedData:     req.AnonymizedData,
	}

	if err := s.userReportRepo.Create(ctx, report); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد گزارش کاربر: %w", err)
	}

	return s.toUserReportResponse(report), nil
}

func (s *UserReportService) GetUserReportById(ctx context.Context, id string) (*schemas.UserReportResponse, error) {
	report, err := s.userReportRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("گزارش کاربر مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت گزارش کاربر: %w", err)
	}

	return s.toUserReportResponse(report), nil
}

func (s *UserReportService) GetAllUserReports(ctx context.Context, req *schemas.UserReportListRequest) (*schemas.UserReportListResponse, error) {
	filterFunc := s.buildUserReportFilters(req)

	reports, total, err := s.userReportRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت گزارش‌های کاربر: %w", err)
	}

	responses := make([]schemas.UserReportResponse, len(reports))
	for i, report := range reports {
		responses[i] = *s.toUserReportResponse(&report)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.UserReportListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.UserReportResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *UserReportService) DeleteUserReport(ctx context.Context, id string) error {
	_, err := s.userReportRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("گزارش کاربر مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت گزارش کاربر: %w", err)
	}

	if err := s.userReportRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف گزارش کاربر: %w", err)
	}

	return nil
}

func (s *UserReportService) buildUserReportFilters(req *schemas.UserReportListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.ReportType != nil {
			db = db.Where("report_type = ?", *req.ReportType)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *UserReportService) toUserReportResponse(report *models.UserReport) *schemas.UserReportResponse {
	return &schemas.UserReportResponse{
		ID:                 report.ID.String(),
		ReportType:         report.ReportType,
		ReportDate:         report.ReportDate,
		TotalUsers:         report.TotalUsers,
		ActiveUsers:        report.ActiveUsers,
		AvgEngagementScore: report.AvgEngagementScore,
		CrisisAlertsCount:  report.CrisisAlertsCount,
		AnonymizedData:     report.AnonymizedData,
		CreatedAt:          report.CreatedAt,
		UpdatedAt:          report.UpdatedAt,
	}
}
