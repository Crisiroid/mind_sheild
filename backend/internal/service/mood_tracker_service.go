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

type MoodTrackerService struct {
	moodTrackerRepo *repository.MoodTrackerRepository
}

func NewMoodTrackerService(moodTrackerRepo *repository.MoodTrackerRepository) *MoodTrackerService {
	return &MoodTrackerService{
		moodTrackerRepo: moodTrackerRepo,
	}
}

func (s *MoodTrackerService) CreateMoodTracker(ctx context.Context, req *schemas.MoodTrackerCreateRequest) (*schemas.MoodTrackerResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	tracker := &models.MoodTracker{
		UserID:       req.UserID,
		ActivityID:   req.ActivityID,
		ActivityName: req.ActivityName,
		MoodBefore:   req.MoodBefore,
		MoodAfter:    req.MoodAfter,
		Notes:        req.Notes,
		DayNumber:    req.DayNumber,
	}

	if err := s.moodTrackerRepo.Create(ctx, tracker); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد ردیاب خلق: %w", err)
	}

	return s.toMoodTrackerResponse(tracker), nil
}

func (s *MoodTrackerService) GetMoodTrackerById(ctx context.Context, id string) (*schemas.MoodTrackerResponse, error) {
	tracker, err := s.moodTrackerRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ردیاب خلق مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت ردیاب خلق: %w", err)
	}

	return s.toMoodTrackerResponse(tracker), nil
}

func (s *MoodTrackerService) GetAllMoodTrackers(ctx context.Context, req *schemas.MoodTrackerListRequest) (*schemas.MoodTrackerListResponse, error) {
	filterFunc := s.buildMoodTrackerFilters(req)

	trackers, total, err := s.moodTrackerRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت ردیاب‌های خلق: %w", err)
	}

	responses := make([]schemas.MoodTrackerResponse, len(trackers))
	for i, tracker := range trackers {
		responses[i] = *s.toMoodTrackerResponse(&tracker)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.MoodTrackerListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.MoodTrackerResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *MoodTrackerService) UpdateMoodTracker(ctx context.Context, id string, req *schemas.MoodTrackerUpdateRequest) (*schemas.MoodTrackerResponse, error) {
	tracker, err := s.moodTrackerRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ردیاب خلق مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت ردیاب خلق: %w", err)
	}

	if req.MoodAfter != nil {
		tracker.MoodAfter = *req.MoodAfter
	}

	if req.Notes != nil {
		tracker.Notes = *req.Notes
	}

	if err := s.moodTrackerRepo.Update(ctx, tracker); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی ردیاب خلق: %w", err)
	}

	return s.toMoodTrackerResponse(tracker), nil
}

func (s *MoodTrackerService) DeleteMoodTracker(ctx context.Context, id string) error {
	_, err := s.moodTrackerRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("ردیاب خلق مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت ردیاب خلق: %w", err)
	}

	if err := s.moodTrackerRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف ردیاب خلق: %w", err)
	}

	return nil
}

func (s *MoodTrackerService) buildMoodTrackerFilters(req *schemas.MoodTrackerListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.ActivityName != nil {
			db = db.Where("activity_name LIKE ?", "%"+*req.ActivityName+"%")
		}

		if req.MoodMin != nil {
			db = db.Where("mood_before >= ? OR mood_after >= ?", *req.MoodMin, *req.MoodMin)
		}

		if req.MoodMax != nil {
			db = db.Where("mood_before <= ? OR mood_after <= ?", *req.MoodMax, *req.MoodMax)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *MoodTrackerService) toMoodTrackerResponse(tracker *models.MoodTracker) *schemas.MoodTrackerResponse {
	return &schemas.MoodTrackerResponse{
		ID:           tracker.ID.String(),
		UserID:       tracker.UserID,
		ActivityID:   tracker.ActivityID,
		ActivityName: tracker.ActivityName,
		MoodBefore:   tracker.MoodBefore,
		MoodAfter:    tracker.MoodAfter,
		MoodDelta:    tracker.MoodAfter - tracker.MoodBefore,
		ActivityDate: tracker.ActivityDate,
		Notes:        tracker.Notes,
		DayNumber:    tracker.DayNumber,
		CreatedAt:    tracker.CreatedAt,
		UpdatedAt:    tracker.UpdatedAt,
	}
}
