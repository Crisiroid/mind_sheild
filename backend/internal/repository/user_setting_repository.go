package repository

import (
	"context"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type UserSettingRepository struct {
	*BaseRepository
}

func NewUserSettingRepository(db *gorm.DB) *UserSettingRepository {
	return &UserSettingRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *UserSettingRepository) Create(ctx context.Context, setting *models.UserSetting) error {
	return r.BaseRepository.Create(ctx, setting)
}

func (r *UserSettingRepository) GetByID(ctx context.Context, id string) (*models.UserSetting, error) {
	var setting models.UserSetting
	err := r.BaseRepository.GetByID(ctx, &setting, id)
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *UserSettingRepository) GetByUserID(ctx context.Context, userID string) (*models.UserSetting, error) {
	var setting models.UserSetting
	err := r.DB.WithContext(ctx).Where("user_id = ?", userID).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *UserSettingRepository) Update(ctx context.Context, setting *models.UserSetting) error {
	return r.BaseRepository.Update(ctx, setting)
}

func (r *UserSettingRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.UserSetting{}, id)
}

func (r *UserSettingRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.UserSetting, int64, error) {
	var settings []models.UserSetting
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.UserSetting{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&settings).Error; err != nil {
		return nil, 0, err
	}

	return settings, total, nil
}

func (r *UserSettingRepository) GetSettingDistribution(ctx context.Context) (map[string][]schemas.DistributionStats, error) {
	result := make(map[string][]schemas.DistributionStats)

	var langDist []schemas.DistributionStats
	r.DB.Model(&models.UserSetting{}).
		Select("language as label, COUNT(*) as count").
		Group("language").
		Order("count DESC").
		Scan(&langDist)
	result["language"] = langDist

	var themeDist []schemas.DistributionStats
	r.DB.Model(&models.UserSetting{}).
		Select("theme as label, COUNT(*) as count").
		Group("theme").
		Order("count DESC").
		Scan(&themeDist)
	result["theme"] = themeDist

	var fontDist []schemas.DistributionStats
	r.DB.Model(&models.UserSetting{}).
		Select("font_size as label, COUNT(*) as count").
		Group("font_size").
		Order("count DESC").
		Scan(&fontDist)
	result["font_size"] = fontDist

	return result, nil
}

func (r *UserSettingRepository) GetNotificationPreferences(ctx context.Context) (*NotificationPreferencesStats, error) {
	var stats NotificationPreferencesStats

	r.DB.Model(&models.UserSetting{}).Where("notification_enabled = ?", true).Count(&stats.NotificationEnabled)
	r.DB.Model(&models.UserSetting{}).Where("notification_enabled = ?", false).Count(&stats.NotificationDisabled)
	r.DB.Model(&models.UserSetting{}).Where("vibration_enabled = ?", true).Count(&stats.VibrationEnabled)
	r.DB.Model(&models.UserSetting{}).Where("vibration_enabled = ?", false).Count(&stats.VibrationDisabled)

	var total int64
	r.DB.Model(&models.UserSetting{}).Count(&total)
	stats.Total = total

	return &stats, nil
}

func (r *UserSettingRepository) GetCrisisAlertSettings(ctx context.Context) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	r.DB.Model(&models.UserSetting{}).
		Select("crisis_alert_threshold as label, COUNT(*) as count").
		Group("crisis_alert_threshold").
		Order("crisis_alert_threshold ASC").
		Scan(&distributions)

	return distributions, nil
}

func (r *UserSettingRepository) GetSettingsByLanguage(ctx context.Context, language string) ([]models.UserSetting, error) {
	var settings []models.UserSetting
	err := r.DB.WithContext(ctx).
		Where("language = ?", language).
		Find(&settings).Error
	return settings, err
}

func (r *UserSettingRepository) GetSettingsByTheme(ctx context.Context, theme string) ([]models.UserSetting, error) {
	var settings []models.UserSetting
	err := r.DB.WithContext(ctx).
		Where("theme = ?", theme).
		Find(&settings).Error
	return settings, err
}

type NotificationPreferencesStats struct {
	Total                int64 `json:"total"`
	NotificationEnabled  int64 `json:"notification_enabled"`
	NotificationDisabled int64 `json:"notification_disabled"`
	VibrationEnabled     int64 `json:"vibration_enabled"`
	VibrationDisabled    int64 `json:"vibration_disabled"`
}
