package repository

import (
	"context"

	"psychology-backend/internal/models"

	"gorm.io/gorm"
)

type AdminRoleRepository struct {
	*BaseRepository
}

func NewAdminRoleRepository(db *gorm.DB) *AdminRoleRepository {
	return &AdminRoleRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *AdminRoleRepository) Create(ctx context.Context, role *models.AdminRole) error {
	return r.BaseRepository.Create(ctx, role)
}

func (r *AdminRoleRepository) GetByID(ctx context.Context, id string) (*models.AdminRole, error) {
	var role models.AdminRole
	err := r.BaseRepository.GetByID(ctx, &role, id)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *AdminRoleRepository) GetByRoleName(ctx context.Context, roleName string) (*models.AdminRole, error) {
	var role models.AdminRole
	err := r.DB.WithContext(ctx).Where("role_name = ?", roleName).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *AdminRoleRepository) Update(ctx context.Context, role *models.AdminRole) error {
	return r.BaseRepository.Update(ctx, role)
}

func (r *AdminRoleRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.AdminRole{}, id)
}

func (r *AdminRoleRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.AdminRole, int64, error) {
	var roles []models.AdminRole
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.AdminRole{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *AdminRoleRepository) GetAllRoles(ctx context.Context) ([]models.AdminRole, error) {
	var roles []models.AdminRole
	err := r.DB.WithContext(ctx).Order("role_name ASC").Find(&roles).Error
	return roles, err
}

func (r *AdminRoleRepository) GetRolePermissions(ctx context.Context) ([]RolePermissionStats, error) {
	var stats []RolePermissionStats

	r.DB.Raw(`
		SELECT
			role_name as label,
			permissions,
			id
		FROM admin_panel.roles
		ORDER BY role_name ASC
	`).Scan(&stats)

	return stats, nil
}

func (r *AdminRoleRepository) GetRoleStatistics(ctx context.Context) ([]RoleAdminCount, error) {
	var stats []RoleAdminCount

	r.DB.Raw(`
		SELECT
			r.role_name as label,
			COUNT(au.id) as count
		FROM admin_panel.roles r
		LEFT JOIN admin_panel.admin_users au ON r.id = au.role_id
		GROUP BY r.id, r.role_name
		ORDER BY count DESC
	`).Scan(&stats)

	return stats, nil
}

type RolePermissionStats struct {
	Label       string `json:"label"`
	Permissions string `json:"permissions"`
	ID          string `json:"id"`
}

type RoleAdminCount struct {
	Label string `json:"label"`
	Count int64  `json:"count"`
}
