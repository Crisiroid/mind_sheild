package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type MindCourtRepository struct {
	*BaseRepository
}

func NewMindCourtRepository(db *gorm.DB) *MindCourtRepository {
	return &MindCourtRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *MindCourtRepository) Create(ctx context.Context, evidence *models.MindCourtEvidence) error {
	return r.BaseRepository.Create(ctx, evidence)
}

func (r *MindCourtRepository) GetByID(ctx context.Context, id string) (*models.MindCourtEvidence, error) {
	var evidence models.MindCourtEvidence
	err := r.BaseRepository.GetByID(ctx, &evidence, id)
	if err != nil {
		return nil, err
	}
	return &evidence, nil
}

func (r *MindCourtRepository) Update(ctx context.Context, evidence *models.MindCourtEvidence) error {
	return r.BaseRepository.Update(ctx, evidence)
}

func (r *MindCourtRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.MindCourtEvidence{}, id)
}

func (r *MindCourtRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.MindCourtEvidence, int64, error) {
	var evidence []models.MindCourtEvidence
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.MindCourtEvidence{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&evidence).Error; err != nil {
		return nil, 0, err
	}

	return evidence, total, nil
}

func (r *MindCourtRepository) GetByUserID(ctx context.Context, userID string) ([]models.MindCourtEvidence, error) {
	var evidence []models.MindCourtEvidence
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&evidence).Error
	return evidence, err
}

func (r *MindCourtRepository) GetEvidenceStats(ctx context.Context, userID string) (*MindCourtStats, error) {
	var stats MindCourtStats

	query := r.DB.Model(&models.MindCourtEvidence{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&stats.TotalEvidence)

	r.DB.Model(&models.MindCourtEvidence{}).
		Where("guide_helper_used = ?", true).
		Count(&stats.GuideHelperUsed)

	if stats.TotalEvidence > 0 {
		stats.GuideHelperRate = float64(stats.GuideHelperUsed) / float64(stats.TotalEvidence) * 100
	}

	return &stats, nil
}

func (r *MindCourtRepository) GetEvidenceTrend(ctx context.Context, userID string, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			COUNT(*) as count
		FROM mind_court_evidence
		WHERE user_id = ? AND created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, userID, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *MindCourtRepository) GetAlternativeThoughtRate(ctx context.Context, userID string) (int64, int64, float64, error) {
	var total, withAlternative int64

	query := r.DB.Model(&models.MindCourtEvidence{})
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)
	r.DB.Model(&models.MindCourtEvidence{}).
		Where("alternative_thought IS NOT NULL AND alternative_thought != ''").
		Count(&withAlternative)

	var rate float64
	if total > 0 {
		rate = float64(withAlternative) / float64(total) * 100
	}

	return total, withAlternative, rate, nil
}

type MindCourtStats struct {
	TotalEvidence   int64   `json:"total_evidence"`
	GuideHelperUsed int64   `json:"guide_helper_used"`
	GuideHelperRate float64 `json:"guide_helper_rate"`
}
