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

type AcceptanceService struct {
	acceptanceRepo interfaces.AcceptanceRepositoryInterface
}

func NewAcceptanceService(acceptanceRepo interfaces.AcceptanceRepositoryInterface) *AcceptanceService {
	return &AcceptanceService{
		acceptanceRepo: acceptanceRepo,
	}
}

func (s *AcceptanceService) CreateAcceptanceExercise(ctx context.Context, req *schemas.AcceptanceExerciseCreateRequest) (*schemas.AcceptanceExerciseResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	exercise := &models.AcceptanceExercise{
		UserID:    req.UserID,
		DayNumber: req.DayNumber,
	}

	if err := s.acceptanceRepo.Create(ctx, exercise); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد تمرین پذیرش: %w", err)
	}

	return s.toAcceptanceExerciseResponse(exercise), nil
}

func (s *AcceptanceService) GetAcceptanceExerciseById(ctx context.Context, id string) (*schemas.AcceptanceExerciseResponse, error) {
	exercise, err := s.acceptanceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تمرین پذیرش مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت تمرین پذیرش: %w", err)
	}

	return s.toAcceptanceExerciseResponse(exercise), nil
}

func (s *AcceptanceService) GetAllAcceptanceExercises(ctx context.Context, req *schemas.AcceptanceExerciseListRequest) (*schemas.AcceptanceExerciseListResponse, error) {
	filterFunc := s.buildAcceptanceExerciseFilters(req)

	exercises, total, err := s.acceptanceRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت تمرین‌های پذیرش: %w", err)
	}

	responses := make([]schemas.AcceptanceExerciseResponse, len(exercises))
	for i, exercise := range exercises {
		responses[i] = *s.toAcceptanceExerciseResponse(&exercise)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.AcceptanceExerciseListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.AcceptanceExerciseResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *AcceptanceService) UpdateAcceptanceExercise(ctx context.Context, id string, req *schemas.AcceptanceExerciseUpdateRequest) (*schemas.AcceptanceExerciseResponse, error) {
	exercise, err := s.acceptanceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تمرین پذیرش مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت تمرین پذیرش: %w", err)
	}

	if req.VideoWatched != nil {
		exercise.VideoWatched = *req.VideoWatched
	}

	if req.UnderstandingLevel != nil {
		exercise.UnderstandingLevel = req.UnderstandingLevel
	}

	if req.Notes != nil {
		exercise.Notes = *req.Notes
	}

	if err := s.acceptanceRepo.Update(ctx, exercise); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی تمرین پذیرش: %w", err)
	}

	return s.toAcceptanceExerciseResponse(exercise), nil
}

func (s *AcceptanceService) DeleteAcceptanceExercise(ctx context.Context, id string) error {
	_, err := s.acceptanceRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("تمرین پذیرش مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت تمرین پذیرش: %w", err)
	}

	if err := s.acceptanceRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف تمرین پذیرش: %w", err)
	}

	return nil
}

func (s *AcceptanceService) buildAcceptanceExerciseFilters(req *schemas.AcceptanceExerciseListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.VideoWatched != nil {
			db = db.Where("video_watched = ?", *req.VideoWatched)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *AcceptanceService) toAcceptanceExerciseResponse(exercise *models.AcceptanceExercise) *schemas.AcceptanceExerciseResponse {
	return &schemas.AcceptanceExerciseResponse{
		ID:                 exercise.ID.String(),
		UserID:             exercise.UserID,
		VideoWatched:       exercise.VideoWatched,
		WatchedAt:          exercise.WatchedAt,
		UnderstandingLevel: exercise.UnderstandingLevel,
		Notes:              exercise.Notes,
		DayNumber:          exercise.DayNumber,
		CreatedAt:          exercise.CreatedAt,
		UpdatedAt:          exercise.UpdatedAt,
	}
}
