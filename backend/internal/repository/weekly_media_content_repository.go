package repository

import (
	"context"
	"fmt"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type WeeklyMediaContentRepository struct {
	*BaseRepository
}

func NewWeeklyMediaContentRepository(db *gorm.DB) *WeeklyMediaContentRepository {
	return &WeeklyMediaContentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *WeeklyMediaContentRepository) Create(ctx context.Context, media *models.WeeklyMediaContent) error {
	return r.BaseRepository.Create(ctx, media)
}

func (r *WeeklyMediaContentRepository) GetByID(ctx context.Context, id string) (*models.WeeklyMediaContent, error) {
	var media models.WeeklyMediaContent
	err := r.BaseRepository.GetByID(ctx, &media, id)
	if err != nil {
		return nil, err
	}
	return &media, nil
}

func (r *WeeklyMediaContentRepository) Update(ctx context.Context, media *models.WeeklyMediaContent) error {
	return r.BaseRepository.Update(ctx, media)
}

func (r *WeeklyMediaContentRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.WeeklyMediaContent{}, id)
}

func (r *WeeklyMediaContentRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.WeeklyMediaContent, int64, error) {
	var mediaList []models.WeeklyMediaContent
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.WeeklyMediaContent{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&mediaList).Error; err != nil {
		return nil, 0, err
	}

	return mediaList, total, nil
}

func (r *WeeklyMediaContentRepository) IncrementDownloadCount(ctx context.Context, id string) error {
	result := r.DB.WithContext(ctx).
		Model(&models.WeeklyMediaContent{}).
		Where("id = ?", id).
		UpdateColumn("download_count", gorm.Expr("download_count + ?", 1))

	if result.Error != nil {
		return fmt.Errorf("خطا در افزایش تعداد دانلود: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *WeeklyMediaContentRepository) GetByWeekNumber(ctx context.Context, weekNumber int, page, pageSize int) ([]models.WeeklyMediaContent, int64, error) {
	return r.List(ctx, page, pageSize, func(db *gorm.DB) *gorm.DB {
		return db.Where("week_number = ? AND is_active = ?", weekNumber, true)
	})
}
