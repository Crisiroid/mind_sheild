package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type WeeklyMediaContentServiceInterface interface {
	CreateMediaContent(ctx context.Context, req *schemas.WeeklyMediaContentCreateRequest, adminID string) (*schemas.WeeklyMediaContentResponse, error)
	UploadMediaContent(ctx context.Context, fileName, fileType, contentType, storagePath string, fileSize int64, req *schemas.WeeklyMediaContentCreateRequest, adminID string) (*schemas.WeeklyMediaContentResponse, error)
	GetMediaContentByID(ctx context.Context, id string) (*schemas.WeeklyMediaContentResponse, error)
	GetAllMediaContent(ctx context.Context, req *schemas.WeeklyMediaContentListRequest) (*schemas.WeeklyMediaContentListResponse, error)
	UpdateMediaContent(ctx context.Context, id string, req *schemas.WeeklyMediaContentUpdateRequest) (*schemas.WeeklyMediaContentResponse, error)
	DeleteMediaContent(ctx context.Context, id string) error
	GetMediaContentByWeek(ctx context.Context, weekNumber int, page, pageSize int) (*schemas.WeeklyMediaContentListResponse, error)
	IncrementDownloadCount(ctx context.Context, id string) error
}
