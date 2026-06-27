package repository

import (
	"context"
	"time"

	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"gorm.io/gorm"
)

type UserRepository struct {
	*BaseRepository
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.BaseRepository.Create(ctx, user)
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := r.BaseRepository.GetByID(ctx, &user, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*models.User, error) {
	var user models.User
	err := r.DB.WithContext(ctx).Where("phone_number = ?", phoneNumber).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByRefreshToken(ctx context.Context, refreshToken string) (*models.User, error) {
	var user models.User
	err := r.DB.WithContext(ctx).
		Where("refresh_token = ?", refreshToken).
		Where("refresh_token_expiry > ?", time.Now()).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	return r.BaseRepository.Update(ctx, user)
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	return r.BaseRepository.Delete(ctx, &models.User{}, id)
}

func (r *UserRepository) List(ctx context.Context, page, pageSize int, filters func(*gorm.DB) *gorm.DB) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := r.DB.WithContext(ctx).Model(&models.User{})

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *UserRepository) GetUserStats(ctx context.Context) (*schemas.UserStatsResponse, error) {
	var stats schemas.UserStatsResponse

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	r.DB.Model(&models.User{}).Count(&stats.TotalUsers)

	thirtyDaysAgo := now.AddDate(0, 0, -30)
	r.DB.Model(&models.User{}).Where("last_login >= ?", thirtyDaysAgo).Count(&stats.ActiveUsers)

	r.DB.Model(&models.User{}).Where("created_at >= ?", today).Count(&stats.NewUsersToday)

	r.DB.Model(&models.User{}).Where("created_at >= ?", weekStart).Count(&stats.NewUsersThisWeek)

	r.DB.Model(&models.User{}).Where("created_at >= ?", monthStart).Count(&stats.NewUsersThisMonth)

	var agreementCount int64
	r.DB.Model(&models.User{}).Where("agreement_accepted = ?", true).Count(&agreementCount)
	if stats.TotalUsers > 0 {
		stats.AgreementRate = float64(agreementCount) / float64(stats.TotalUsers) * 100
	}

	r.DB.Model(&models.User{}).Select("COALESCE(AVG(login_count), 0)").Scan(&stats.AvgLoginCount)

	return &stats, nil
}

func (r *UserRepository) GetUserActivityTrend(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(created_at) as timestamp,
			COUNT(*) as count
		FROM users
		WHERE created_at >= ? AND created_at <= ?
		GROUP BY DATE(created_at)
		ORDER BY timestamp ASC
	`, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *UserRepository) GetLoginAnalytics(ctx context.Context, dateFrom, dateTo time.Time) ([]schemas.TrendDataPoint, error) {
	var trends []schemas.TrendDataPoint

	r.DB.Raw(`
		SELECT
			DATE(last_login) as timestamp,
			COUNT(*) as count
		FROM users
		WHERE last_login >= ? AND last_login <= ?
		GROUP BY DATE(last_login)
		ORDER BY timestamp ASC
	`, dateFrom, dateTo).Scan(&trends)

	return trends, nil
}

func (r *UserRepository) GetAgreementStats(ctx context.Context) (int64, int64, float64, error) {
	var totalUsers, agreedUsers int64

	r.DB.Model(&models.User{}).Count(&totalUsers)
	r.DB.Model(&models.User{}).Where("agreement_accepted = ?", true).Count(&agreedUsers)

	var rate float64
	if totalUsers > 0 {
		rate = float64(agreedUsers) / float64(totalUsers) * 100
	}

	return totalUsers, agreedUsers, rate, nil
}

func (r *UserRepository) GetAppVersionDistribution(ctx context.Context) ([]schemas.DistributionStats, error) {
	var distributions []schemas.DistributionStats

	r.DB.Model(&models.User{}).
		Select("COALESCE(app_version, 'unknown') as label, COUNT(*) as count").
		Group("app_version").
		Order("count DESC").
		Scan(&distributions)

	return distributions, nil
}

func (r *UserRepository) GetInactiveUsers(ctx context.Context, daysThreshold int) ([]models.User, error) {
	var users []models.User
	cutoffDate := time.Now().AddDate(0, 0, -daysThreshold)

	err := r.DB.WithContext(ctx).
		Where("last_login < ? OR last_login IS NULL", cutoffDate).
		Find(&users).Error

	return users, err
}

func (r *UserRepository) ExportUsers(ctx context.Context, dateFrom, dateTo *time.Time, userID string) ([]models.User, error) {
	var users []models.User
	query := r.DB.WithContext(ctx).Model(&models.User{})

	if dateFrom != nil {
		query = query.Where("created_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("created_at <= ?", *dateTo)
	}
	if userID != "" {
		query = query.Where("id = ?", userID)
	}

	err := query.Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *UserRepository) GetUserEngagementStats(ctx context.Context, dateFrom, dateTo time.Time) (*schemas.EngagementStatsResponse, error) {
	var stats schemas.EngagementStatsResponse

	r.DB.Model(&models.User{}).
		Where("created_at >= ? AND created_at <= ?", dateFrom, dateTo).
		Count(&stats.TotalUsers)

	r.DB.Model(&models.User{}).
		Where("last_login >= ? AND last_login <= ?", dateFrom, dateTo).
		Count(&stats.ActiveUsers)

	if stats.TotalUsers > 0 {
		stats.EngagementRate = float64(stats.ActiveUsers) / float64(stats.TotalUsers) * 100
	}

	return &stats, nil
}
