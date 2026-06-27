package repository

import (
	"context"
	"math"
	"time"

	"gorm.io/gorm"
)

type BaseRepository struct {
	DB *gorm.DB
}

func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{DB: db}
}

func (r *BaseRepository) Create(ctx context.Context, entity interface{}) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository) GetByID(ctx context.Context, entity interface{}, id string) error {
	return r.DB.WithContext(ctx).First(entity, "id = ?", id).Error
}

func (r *BaseRepository) Update(ctx context.Context, entity interface{}) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository) Delete(ctx context.Context, entity interface{}, id string) error {
	return r.DB.WithContext(ctx).Delete(entity, "id = ?", id).Error
}

type PaginatedResult struct {
	Data     []interface{}
	Total    int64
	Page     int
	PageSize int
	Pages    int
}

func (r *BaseRepository) ListWithPagination(ctx context.Context, entity interface{}, page, pageSize int, filters func(*gorm.DB) *gorm.DB) (*PaginatedResult, error) {
	var total int64
	var data []interface{}

	query := r.DB.WithContext(ctx).Model(entity)

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&data).Error; err != nil {
		return nil, err
	}

	pages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &PaginatedResult{
		Data:     data,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}, nil
}

func (r *BaseRepository) Count(ctx context.Context, entity interface{}, filters func(*gorm.DB) *gorm.DB) (int64, error) {
	var total int64
	query := r.DB.WithContext(ctx).Model(entity)

	if filters != nil {
		query = filters(query)
	}

	err := query.Count(&total).Error
	return total, err
}

func ApplyUserIDFilter(db *gorm.DB, userID string) *gorm.DB {
	if userID != "" {
		return db.Where("user_id = ?", userID)
	}
	return db
}

func ApplyDateRangeFilter(db *gorm.DB, dateFrom, dateTo *time.Time) *gorm.DB {
	if dateFrom != nil {
		db = db.Where("created_at >= ?", *dateFrom)
	}
	if dateTo != nil {
		db = db.Where("created_at <= ?", *dateTo)
	}
	return db
}

func ApplySearchFilter(db *gorm.DB, search string, fields ...string) *gorm.DB {
	if search != "" && len(fields) > 0 {
		searchPattern := "%" + search + "%"
		var conditions []string
		var values []interface{}

		for _, field := range fields {
			conditions = append(conditions, field+" LIKE ?")
			values = append(values, searchPattern)
		}

		query := "(" + joinOr(conditions) + ")"
		db = db.Where(query, values...)
	}
	return db
}

func ApplyDayNumberFilter(db *gorm.DB, dayNumber *int) *gorm.DB {
	if dayNumber != nil {
		return db.Where("day_number = ?", *dayNumber)
	}
	return db
}

func ApplyBoolFilter(db *gorm.DB, field string, value *bool) *gorm.DB {
	if value != nil {
		return db.Where(field+" = ?", *value)
	}
	return db
}

func ApplyIntRangeFilter(db *gorm.DB, field string, min, max *int) *gorm.DB {
	if min != nil {
		db = db.Where(field+" >= ?", *min)
	}
	if max != nil {
		db = db.Where(field+" <= ?", *max)
	}
	return db
}

func ApplySorting(db *gorm.DB, sortBy, sortOrder string) *gorm.DB {
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}
	return db.Order(sortBy + " " + sortOrder)
}

func ApplyPagination(db *gorm.DB, page, pageSize int) *gorm.DB {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return db.Offset(offset).Limit(pageSize)
}

func (r *BaseRepository) GetWithPreloads(ctx context.Context, entity interface{}, id string, preloads ...string) error {
	query := r.DB.WithContext(ctx)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return query.First(entity, "id = ?", id).Error
}

func (r *BaseRepository) ListWithPreloads(ctx context.Context, entity interface{}, page, pageSize int, preloads []string, filters func(*gorm.DB) *gorm.DB) (*PaginatedResult, error) {
	var total int64

	query := r.DB.WithContext(ctx).Model(entity)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if filters != nil {
		query = filters(query)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var data []interface{}
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&data).Error; err != nil {
		return nil, err
	}

	pages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &PaginatedResult{
		Data:     data,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		Pages:    pages,
	}, nil
}

func (r *BaseRepository) BatchCreate(ctx context.Context, entities interface{}) error {
	return r.DB.WithContext(ctx).Create(entities).Error
}

func (r *BaseRepository) UpdateField(ctx context.Context, entity interface{}, id string, field string, value interface{}) error {
	return r.DB.WithContext(ctx).Model(entity).Where("id = ?", id).Update(field, value).Error
}

func (r *BaseRepository) Increment(ctx context.Context, entity interface{}, id string, field string, value int64) error {
	return r.DB.WithContext(ctx).Model(entity).Where("id = ?", id).Update(field, gorm.Expr(field+" + ?", value)).Error
}

func (r *BaseRepository) Exists(ctx context.Context, entity interface{}, conditions string, args ...interface{}) (bool, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(entity).Where(conditions, args...).Count(&count).Error
	return count > 0, err
}

func joinOr(conditions []string) string {
	result := ""
	for i, cond := range conditions {
		if i > 0 {
			result += " OR "
		}
		result += cond
	}
	return result
}

func (r *BaseRepository) FirstOrCreate(ctx context.Context, entity interface{}, where interface{}) error {
	return r.DB.WithContext(ctx).FirstOrCreate(entity, where).Error
}

func (r *BaseRepository) UpdateOrCreate(ctx context.Context, entity interface{}, where interface{}, data interface{}) error {
	return r.DB.WithContext(ctx).Where(where).Assign(data).FirstOrCreate(entity).Error
}

func (r *BaseRepository) GroupBy(ctx context.Context, entity interface{}, groupBy string, selectClause string, dest interface{}) error {
	return r.DB.WithContext(ctx).Model(entity).Select(selectClause).Group(groupBy).Find(dest).Error
}

func (r *BaseRepository) Raw(ctx context.Context, query string, dest interface{}, args ...interface{}) error {
	return r.DB.WithContext(ctx).Raw(query, args...).Scan(dest).Error
}
