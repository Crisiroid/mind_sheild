package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"psychology-backend/internal/models"
	"psychology-backend/internal/repository"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type EmotionTriangleService struct {
	emotionRepo *repository.EmotionTriangleRepository
}

func NewEmotionTriangleService(emotionRepo *repository.EmotionTriangleRepository) *EmotionTriangleService {
	return &EmotionTriangleService{
		emotionRepo: emotionRepo,
	}
}

func (s *EmotionTriangleService) CreateEmotionInteraction(ctx context.Context, req *schemas.EmotionInteractionCreateRequest) (*schemas.EmotionInteractionResponse, error) {
	if req.UserID == "" {
		return nil, errors.New("شناسه کاربر الزامی است")
	}

	if req.SideClicked == "" {
		return nil, errors.New("ضلع انتخاب شده الزامی است")
	}

	thoughtAccountsJSON := ""
	if len(req.ThoughtAccountsViewed) > 0 {
		if data, err := json.Marshal(req.ThoughtAccountsViewed); err == nil {
			thoughtAccountsJSON = string(data)
		}
	}

	interaction := &models.EmotionTriangleInteraction{
		UserID:                req.UserID,
		SideClicked:           req.SideClicked,
		ThoughtAccountsViewed: thoughtAccountsJSON,
		VibrationTriggered:    req.VibrationTriggered,
		DayNumber:             req.DayNumber,
	}

	if err := s.emotionRepo.Create(ctx, interaction); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد تعامل احساسی: %w", err)
	}

	return s.toEmotionInteractionResponse(interaction), nil
}

func (s *EmotionTriangleService) GetEmotionInteractionById(ctx context.Context, id string) (*schemas.EmotionInteractionResponse, error) {
	interaction, err := s.emotionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("تعامل احساسی مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت تعامل احساسی: %w", err)
	}

	return s.toEmotionInteractionResponse(interaction), nil
}

func (s *EmotionTriangleService) GetAllEmotionInteractions(ctx context.Context, req *schemas.EmotionInteractionListRequest) (*schemas.EmotionInteractionListResponse, error) {
	filterFunc := s.buildEmotionInteractionFilters(req)

	interactions, total, err := s.emotionRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت تعاملات احساسی: %w", err)
	}

	responses := make([]schemas.EmotionInteractionResponse, len(interactions))
	for i, interaction := range interactions {
		responses[i] = *s.toEmotionInteractionResponse(&interaction)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.EmotionInteractionListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.EmotionInteractionResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *EmotionTriangleService) DeleteEmotionInteraction(ctx context.Context, id string) error {
	_, err := s.emotionRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("تعامل احساسی مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت تعامل احساسی: %w", err)
	}

	if err := s.emotionRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف تعامل احساسی: %w", err)
	}

	return nil
}

func (s *EmotionTriangleService) buildEmotionInteractionFilters(req *schemas.EmotionInteractionListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = repository.ApplyDateRangeFilter(db, req.DateFrom, req.DateTo)

		if req.SideClicked != nil {
			db = db.Where("side_clicked = ?", *req.SideClicked)
		}

		if req.SortBy != "" {
			db = repository.ApplySorting(db, req.SortBy, req.SortOrder)
		} else {
			db = db.Order("created_at desc")
		}

		return db
	}
}

func (s *EmotionTriangleService) toEmotionInteractionResponse(interaction *models.EmotionTriangleInteraction) *schemas.EmotionInteractionResponse {
	var thoughtAccounts []string
	if interaction.ThoughtAccountsViewed != "" {
		json.Unmarshal([]byte(interaction.ThoughtAccountsViewed), &thoughtAccounts)
	}

	return &schemas.EmotionInteractionResponse{
		ID:                    interaction.ID.String(),
		UserID:                interaction.UserID,
		InteractionDate:       interaction.InteractionDate,
		SideClicked:           interaction.SideClicked,
		ThoughtAccountsViewed: thoughtAccounts,
		VibrationTriggered:    interaction.VibrationTriggered,
		DayNumber:             interaction.DayNumber,
		CreatedAt:             interaction.CreatedAt,
		UpdatedAt:             interaction.UpdatedAt,
	}
}
