package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type AdminUserRepository struct {
	*BaseRepository
}

func NewAdminUserRepository(db *gorm.DB) *AdminUserRepository {
	return &AdminUserRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *AdminUserRepository) Create(ctx context.Context, adminUser *models.AdminUser) error {
	return r.BaseRepository.Create(ctx, adminUser)
}

func (r *AdminUserRepository) GetByID(ctx context.Context, id string) (*models.AdminUser, error) {
	var adminUser models.AdminUser
	err := r.BaseRepository.GetByID(ctx, &adminUser, id)
	if err != nil {
		return nil, err
	}
	return &adminUser, nil
}

func (r *AdminUserRepository) GetByUsername(ctx context.Context, username string) (*models.AdminUser, error) {
	var adminUser models.AdminUser
	err := r.DB.WithContext(ctx).Where("username = ?", username).First(&adminUser).Error
	if err != nil {
		return nil, err
	}
	return &adminUser, nil
}

func (r *AdminUserRepository) GetByEmail(ctx context.Context, email string) (*models.AdminUser, error) {
	var adminUser models.AdminUser
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&adminUser).Error
	if err != nil {
		return nil, err
	}
	return &adminUser, nil
}

func (r *AdminUserRepository) Update(ctx context.Context, adminUser *models.AdminUser) error {
	return r.BaseRepository.Update(ctx, adminUser)
}

func (r *AdminUserRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.AdminUser{}, id)
}

func (r *AdminUserRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.AdminUser, int64, error) {
	var adminUsers []models.AdminUser
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.AdminUser{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&adminUsers).Error; err != nil {
		return nil, 0, err
	}

	return adminUsers, total, nil
}

func (r *AdminUserRepository) GetAdminStats(ctx context.Context) (int64, int64, error) {
	var totalAdmins, activeAdmins int64

	r.DB.Model(&models.AdminUser{}).Count(&totalAdmins)
	r.DB.Model(&models.AdminUser{}).Where("is_active = ?", true).Count(&activeAdmins)

	return totalAdmins, activeAdmins, nil
}

func (r *AdminUserRepository) GetAdminActivityLog(ctx context.Context, dateFrom, dateTo time.Time) ([]AdminActivityLog, error) {
	var logs []AdminActivityLog

	r.DB.Raw(`
		SELECT
			last_login as timestamp,
			COUNT(*) as count
		FROM admin_panel.admin_users
		WHERE last_login >= ? AND last_login <= ? AND last_login != ''
		GROUP BY last_login
		ORDER BY timestamp ASC
	`, dateFrom, dateTo).Scan(&logs)

	return logs, nil
}

func (r *AdminUserRepository) GetAdminsByRole(ctx context.Context) ([]RoleDistribution, error) {
	var distributions []RoleDistribution

	r.DB.Raw(`
		SELECT
			role_id as label,
			COUNT(*) as count
		FROM admin_panel.admin_users
		GROUP BY role_id
		ORDER BY count DESC
	`).Scan(&distributions)

	return distributions, nil
}

func (r *AdminUserRepository) GetInactiveAdmins(ctx context.Context, daysThreshold int) ([]models.AdminUser, error) {
	var adminUsers []models.AdminUser
	err := r.DB.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&adminUsers).Error

	return adminUsers, err
}

func (r *AdminUserRepository) GetAdminsByActiveStatus(ctx context.Context, isActive bool) ([]models.AdminUser, error) {
	var adminUsers []models.AdminUser
	err := r.DB.WithContext(ctx).
		Where("is_active = ?", isActive).
		Order("created_at DESC").
		Find(&adminUsers).Error

	return adminUsers, err
}

type AdminActivityLog struct {
	Timestamp string `json:"timestamp"`
	Count     int64  `json:"count"`
}

type RoleDistribution struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}
