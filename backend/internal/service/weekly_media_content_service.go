package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"psychology-backend/internal/interfaces"
	"psychology-backend/internal/models"
	"psychology-backend/pkg/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WeeklyMediaContentService struct {
	mediaRepo interfaces.WeeklyMediaContentRepositoryInterface
}

func NewWeeklyMediaContentService(mediaRepo interfaces.WeeklyMediaContentRepositoryInterface) *WeeklyMediaContentService {
	return &WeeklyMediaContentService{
		mediaRepo: mediaRepo,
	}
}

func (s *WeeklyMediaContentService) CreateMediaContent(ctx context.Context, req *schemas.WeeklyMediaContentCreateRequest, adminID string) (*schemas.WeeklyMediaContentResponse, error) {
	return nil, errors.New("استفاده از UploadMediaContent برای آپلود فایل")
}

func (s *WeeklyMediaContentService) UploadMediaContent(ctx context.Context, fileName, fileType, contentType, storagePath string, fileSize int64, req *schemas.WeeklyMediaContentCreateRequest, adminID string) (*schemas.WeeklyMediaContentResponse, error) {
	if adminID == "" {
		return nil, errors.New("شناسه ادمین الزامی است")
	}

	ext := filepath.Ext(fileName)
	nameWithoutExt := strings.TrimSuffix(fileName, ext)
	uniqueFileName := fmt.Sprintf("%s_%s%s", nameWithoutExt, uuid.New().String()[:8], ext)

	fileURL := fmt.Sprintf("/api/v1/media/weekly/files/%s", uniqueFileName)

	media := &models.WeeklyMediaContent{
		FileName:       uniqueFileName,
		FileType:       fileType,
		FileSize:       fileSize,
		FileURL:        fileURL,
		WeekNumber:     req.WeekNumber,
		Description:    req.Description,
		ContentType:    contentType,
		OriginalName:   fileName,
		StoragePath:    storagePath,
		IsActive:       true,
		DownloadCount:  0,
		CreatedByAdmin: adminID,
	}

	if err := s.mediaRepo.Create(ctx, media); err != nil {
		return nil, fmt.Errorf("خطا در ایجاد محتوای مدیا: %w", err)
	}

	return s.toMediaResponse(media), nil
}

func (s *WeeklyMediaContentService) GetMediaContentByID(ctx context.Context, id string) (*schemas.WeeklyMediaContentResponse, error) {
	media, err := s.mediaRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("محتوای مدیا مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت محتوای مدیا: %w", err)
	}

	return s.toMediaResponse(media), nil
}

func (s *WeeklyMediaContentService) GetAllMediaContent(ctx context.Context, req *schemas.WeeklyMediaContentListRequest) (*schemas.WeeklyMediaContentListResponse, error) {
	filterFunc := s.buildMediaFilters(req)

	mediaList, total, err := s.mediaRepo.List(ctx, req.Page, req.PageSize, filterFunc)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت محتوای مدیا: %w", err)
	}

	responses := make([]schemas.WeeklyMediaContentResponse, len(mediaList))
	for i, media := range mediaList {
		responses[i] = *s.toMediaResponse(&media)
	}

	pages := int((total + int64(req.PageSize) - 1) / int64(req.PageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.WeeklyMediaContentListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.WeeklyMediaContentResponse]{
			Data:     responses,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *WeeklyMediaContentService) UpdateMediaContent(ctx context.Context, id string, req *schemas.WeeklyMediaContentUpdateRequest) (*schemas.WeeklyMediaContentResponse, error) {
	media, err := s.mediaRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("محتوای مدیا مورد نظر یافت نشد")
		}
		return nil, fmt.Errorf("خطا در دریافت محتوای مدیا: %w", err)
	}

	if req.WeekNumber > 0 {
		media.WeekNumber = req.WeekNumber
	}
	if req.FileType != "" {
		media.FileType = req.FileType
	}
	if req.Description != "" {
		media.Description = req.Description
	}
	if req.IsActive != nil {
		media.IsActive = *req.IsActive
	}

	if err := s.mediaRepo.Update(ctx, media); err != nil {
		return nil, fmt.Errorf("خطا در بروزرسانی محتوای مدیا: %w", err)
	}

	return s.toMediaResponse(media), nil
}

func (s *WeeklyMediaContentService) DeleteMediaContent(ctx context.Context, id string) error {
	_, err := s.mediaRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("محتوای مدیا مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت محتوای مدیا: %w", err)
	}

	if err := s.mediaRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("خطا در حذف محتوای مدیا: %w", err)
	}

	return nil
}

func (s *WeeklyMediaContentService) GetMediaContentByWeek(ctx context.Context, weekNumber int, page, pageSize int) (*schemas.WeeklyMediaContentListResponse, error) {
	mediaList, total, err := s.mediaRepo.GetByWeekNumber(ctx, weekNumber, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("خطا در دریافت محتوای مدیا: %w", err)
	}

	responses := make([]schemas.WeeklyMediaContentResponse, len(mediaList))
	for i, media := range mediaList {
		responses[i] = *s.toMediaResponse(&media)
	}

	pages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if pages == 0 {
		pages = 1
	}

	return &schemas.WeeklyMediaContentListResponse{
		PaginatedResponse: schemas.PaginatedResponse[schemas.WeeklyMediaContentResponse]{
			Data:     responses,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
			Pages:    pages,
		},
	}, nil
}

func (s *WeeklyMediaContentService) IncrementDownloadCount(ctx context.Context, id string) error {
	media, err := s.mediaRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("محتوای مدیا مورد نظر یافت نشد")
		}
		return fmt.Errorf("خطا در دریافت محتوای مدیا: %w", err)
	}

	if !media.IsActive {
		return errors.New("این محتوا غیرفعال است")
	}

	if err := s.mediaRepo.IncrementDownloadCount(ctx, id); err != nil {
		return fmt.Errorf("خطا در افزایش تعداد دانلود: %w", err)
	}

	return nil
}

func (s *WeeklyMediaContentService) buildMediaFilters(req *schemas.WeeklyMediaContentListRequest) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if req.WeekNumber != nil {
			db = db.Where("week_number = ?", *req.WeekNumber)
		}
		if req.FileType != "" {
			db = db.Where("file_type = ?", req.FileType)
		}
		if req.IsActive != nil {
			db = db.Where("is_active = ?", *req.IsActive)
		}
		return db
	}
}

func (s *WeeklyMediaContentService) toMediaResponse(media *models.WeeklyMediaContent) *schemas.WeeklyMediaContentResponse {
	return &schemas.WeeklyMediaContentResponse{
		ID:            media.ID.String(),
		FileName:      media.FileName,
		FileType:      media.FileType,
		FileSize:      media.FileSize,
		FileURL:       media.FileURL,
		WeekNumber:    media.WeekNumber,
		Description:   media.Description,
		ContentType:   media.ContentType,
		OriginalName:  media.OriginalName,
		IsActive:      media.IsActive,
		DownloadCount: media.DownloadCount,
		CreatedAt:     media.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     media.UpdatedAt.Format(time.RFC3339),
	}
}
