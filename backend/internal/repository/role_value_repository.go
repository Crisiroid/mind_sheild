package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type RoleValueRepository struct {
	*BaseRepository
}

func NewRoleValueRepository(db *gorm.DB) *RoleValueRepository {
	return &RoleValueRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *RoleValueRepository) Create(ctx context.Context, roleValue *models.RoleAndValue) error {
	return r.BaseRepository.Create(ctx, roleValue)
}

func (r *RoleValueRepository) GetByID(ctx context.Context, id string) (*models.RoleAndValue, error) {
	var roleValue models.RoleAndValue
	err := r.BaseRepository.GetByID(ctx, &roleValue, id)
	if err != nil {
		return nil, err
	}
	return &roleValue, nil
}

func (r *RoleValueRepository) Update(ctx context.Context, roleValue *models.RoleAndValue) error {
	return r.BaseRepository.Update(ctx, roleValue)
}

func (r *RoleValueRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.RoleAndValue{}, id)
}

func (r *RoleValueRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.RoleAndValue, int64, error) {
	var roleValues []models.RoleAndValue
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.RoleAndValue{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&roleValues).Error; err != nil {
		return nil, 0, err
	}

	return roleValues, total, nil
}

func (r *RoleValueRepository) GetByUserID(ctx context.Context, userID string) ([]models.RoleAndValue, error) {
	var roleValues []models.RoleAndValue
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&roleValues).Error
	return roleValues, err
}

func (r *RoleValueRepository) GetRoleValueStats(ctx context.Context, userID string) (*RoleValueStats, error) {
	var stats RoleValueStats

	query := r.DB.Model(&models.RoleAndValue{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&stats.TotalEntries)

	r.DB.Model(&models.RoleAndValue{}).
		Where("entry_type = ?", "role").
		Count(&stats.RoleCount)

	r.DB.Model(&models.RoleAndValue{}).
		Where("entry_type = ?", "value").
		Count(&stats.ValueCount)

	return &stats, nil
}

func (r *RoleValueRepository) GetEntryTypeDistribution(ctx context.Context, userID string) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	query := r.DB.Model(&models.RoleAndValue{}).
		Select("entry_type as label, COUNT(*) as count")

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Group("entry_type").
		Order("count DESC").
		Scan(&distributions)

	return distributions, nil
}

func (r *RoleValueRepository) GetEntriesByType(ctx context.Context, userID, entryType string) ([]models.RoleAndValue, error) {
	var roleValues []models.RoleAndValue
	query := r.DB.WithContext(ctx).Where("entry_type = ?", entryType)

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	err := query.Order("created_at DESC").Find(&roleValues).Error
	return roleValues, err
}

func (r *RoleValueRepository) GetEntryTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			COUNT(*) as count
		FROM roles_and_values
		WHERE user_id = ? AND created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, userID, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

type RoleValueStats struct {
	TotalEntries int64 `json:"total_entries"`
	RoleCount    int64 `json:"role_count"`
	ValueCount   int64 `json:"value_count"`
}
