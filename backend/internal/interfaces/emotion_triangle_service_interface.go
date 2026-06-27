package interfaces

import (
	"context"

	"psychology-backend/pkg/schemas"
)

type EmotionTriangleServiceInterface interface {
	CreateEmotionInteraction(ctx context.Context, req *schemas.EmotionInteractionCreateRequest) (*schemas.EmotionInteractionResponse, error)
	GetEmotionInteractionById(ctx context.Context, id string) (*schemas.EmotionInteractionResponse, error)
	GetAllEmotionInteractions(ctx context.Context, req *schemas.EmotionInteractionListRequest) (*schemas.EmotionInteractionListResponse, error)
	UpdateEmotionInteraction(ctx context.Context, id string, req *schemas.EmotionInteractionUpdateRequest) (*schemas.EmotionInteractionResponse, error)
	DeleteEmotionInteraction(ctx context.Context, id string) error
}
