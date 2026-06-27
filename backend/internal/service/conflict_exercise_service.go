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

type ConflictExerciseService struct {
	conflictRepo interfaces.ConflictExerciseRepositoryInterface
}

func NewConflictExerciseService(conflictRepo interfaces.ConflictExerciseRepositoryInterface) *ConflictExerciseService {
	return &ConflictExerciseService{
		conflictRepo: conflictRepo,
	}
}

func (s *ConflictExerciseService) CreateConflictExercise(ctx context.Context, req *schemas.ConflictExerciseCreateRequest) (*schemas.ConflictExerciseResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	exercise := &models.ConflictExercise{
		UserID:           req.UserID,
		ScenarioID:       req.ScenarioID,
		PerformanceScore: req.PerformanceScore,
		DayNumber:        req.DayNumber,
	}

	if err := s.conflictRepo.Create(ctx, exercise); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد تمرین تعارض: %w", err)
	}

	return s.toConflictExerciseResponse(exercise), nil
}

func (s *ConflictExerciseService) GetConflictExerciseById(ctx context.Context, id string) (*schemas.ConflictExerciseResponse, error) {
	exercise, err := s.conflictRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تمرین تعارض مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت تمرین تعارض: %w", err)
	}

	return s.toConflictExerciseResponse(exercise), nil
}

func (s *ConflictExerciseService) GetAllConflictExercises(ctx context.Context, req *schemas.ConflictExerciseListRequest) (*schemas.ConflictExerciseListResponse, error) {
	filterFunc := s.buildConflictExerciseFilters(req)

	exercises, total, err := s.conflictRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت تمرین‌های تعارض: %w", err)
	}

	responses := make([]schemas.ConflictExerciseResponse, len(exercises))
	for i, exercise := range exercises {
		responses[i] = *s.toConflictExerciseResponse(&exercise)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.ConflictExerciseListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.ConflictExerciseResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *ConflictExerciseService) UpdateConflictExercise(ctx context.Context, id string, req *schemas.ConflictExerciseUpdateRequest) (*schemas.ConflictExerciseResponse, error) {
	exercise, err := s.conflictRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تمرین تعارض مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت تمرین تعارض: %w", err)
	}

	if req.PracticeCount != nil {
		exercise.PracticeCount = *req.PracticeCount
	}

	if req.PerformanceScore != nil {
		exercise.PerformanceScore = req.PerformanceScore
	}

	if err := s.conflictRepo.Update(ctx, exercise); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی تمرین تعارض: %w", err)
	}

	return s.toConflictExerciseResponse(exercise), nil
}

func (s *ConflictExerciseService) DeleteConflictExercise(ctx context.Context, id string) error {
	_, err := s.conflictRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("تمرین تعارض مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت تمرین تعارض: %w", err)
	}

	if err := s.conflictRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف تمرین تعارض: %w", err)
	}

	return nil
}

func (s *ConflictExerciseService) buildConflictExerciseFilters(req *schemas.ConflictExerciseListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.ScenarioID != nil {
			db = db.Where("scenario_id = ?", *req.ScenarioID)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *ConflictExerciseService) toConflictExerciseResponse(exercise *models.ConflictExercise) *schemas.ConflictExerciseResponse {
	return &schemas.ConflictExerciseResponse{
		ID:               exercise.ID.String(),
		UserID:           exercise.UserID,
		ScenarioID:       exercise.ScenarioID,
		PracticeCount:    exercise.PracticeCount,
		LastPracticeDate: exercise.LastPracticeDate,
		PerformanceScore: exercise.PerformanceScore,
		DayNumber:        exercise.DayNumber,
		CreatedAt:        exercise.CreatedAt,
		UpdatedAt:        exercise.UpdatedAt,
	}
}
