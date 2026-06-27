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

type MindCourtService struct {
	mindCourtRepo *repository.MindCourtRepository
}

func NewMindCourtService(mindCourtRepo *repository.MindCourtRepository) *MindCourtService {
	return &MindCourtService{
		mindCourtRepo: mindCourtRepo,
	}
}

func (s *MindCourtService) CreateMindCourt(ctx context.Context, req *schemas.MindCourtCreateRequest) (*schemas.MindCourtResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	if req.NegativeThoughtID == "" {
		return nil, errors.New("شناسه فکر منفی الزامی است")
	}

	mindCourt := &models.MindCourtEvidence{
		UserID:                req.UserID,
		NegativeThoughtID:     req.NegativeThoughtID,
		SupportingEvidence:    req.SupportingEvidence,
		ContradictingEvidence: req.ContradictingEvidence,
		GuideHelperUsed:       req.GuideHelperUsed,
		AlternativeThought:    req.AlternativeThought,
		DayNumber:             req.DayNumber,
	}

	if err := s.mindCourtRepo.Create(ctx, mindCourt); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد دادگاه ذهن: %w", err)
	}

	return s.toMindCourtResponse(mindCourt), nil
}

func (s *MindCourtService) GetMindCourtById(ctx context.Context, id string) (*schemas.MindCourtResponse, error) {
	mindCourt, err := s.mindCourtRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("دادگاه ذهن مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت دادگاه ذهن: %w", err)
	}

	return s.toMindCourtResponse(mindCourt), nil
}

func (s *MindCourtService) GetAllMindCourts(ctx context.Context, req *schemas.MindCourtListRequest) (*schemas.MindCourtListResponse, error) {
	filterFunc := s.buildMindCourtFilters(req)

	mindCourts, total, err := s.mindCourtRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت دادگاه‌های ذهن: %w", err)
	}

	responses := make([]schemas.MindCourtResponse, len(mindCourts))
	for i, mindCourt := range mindCourts {
		responses[i] = *s.toMindCourtResponse(&mindCourt)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.MindCourtListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.MindCourtResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *MindCourtService) UpdateMindCourt(ctx context.Context, id string, req *schemas.MindCourtUpdateRequest) (*schemas.MindCourtResponse, error) {
	mindCourt, err := s.mindCourtRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("دادگاه ذهن مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت دادگاه ذهن: %w", err)
	}

	if req.SupportingEvidence != nil {
		mindCourt.SupportingEvidence = *req.SupportingEvidence
	}

	if req.ContradictingEvidence != nil {
		mindCourt.ContradictingEvidence = *req.ContradictingEvidence
	}

	if req.AlternativeThought != nil {
		mindCourt.AlternativeThought = *req.AlternativeThought
	}

	if err := s.mindCourtRepo.Update(ctx, mindCourt); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی دادگاه ذهن: %w", err)
	}

	return s.toMindCourtResponse(mindCourt), nil
}

func (s *MindCourtService) DeleteMindCourt(ctx context.Context, id string) error {
	_, err := s.mindCourtRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("دادگاه ذهن مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت دادگاه ذهن: %w", err)
	}

	if err := s.mindCourtRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف دادگاه ذهن: %w", err)
	}

	return nil
}

func (s *MindCourtService) buildMindCourtFilters(req *schemas.MindCourtListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.NegativeThoughtID != nil {
			db = db.Where("negative_thought_id = ?", *req.NegativeThoughtID)
		}

		if req.GuideHelperUsed != nil {
			db = db.Where("guide_helper_used = ?", *req.GuideHelperUsed)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *MindCourtService) toMindCourtResponse(mindCourt *models.MindCourtEvidence) *schemas.MindCourtResponse {
	return &schemas.MindCourtResponse{
		ID:                    mindCourt.ID.String(),
		UserID:                mindCourt.UserID,
		NegativeThoughtID:     mindCourt.NegativeThoughtID,
		SupportingEvidence:    mindCourt.SupportingEvidence,
		ContradictingEvidence: mindCourt.ContradictingEvidence,
		GuideHelperUsed:       mindCourt.GuideHelperUsed,
		AlternativeThought:    mindCourt.AlternativeThought,
		CreatedDate:           mindCourt.CreatedDate,
		DayNumber:             mindCourt.DayNumber,
		CreatedAt:             mindCourt.CreatedAt,
		UpdatedAt:             mindCourt.UpdatedAt,
	}
}
