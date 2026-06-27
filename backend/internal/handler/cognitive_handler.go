package handler

import (
	"net/http"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type CognitiveHandler struct {
	*BaseHandler
	cognitiveService interfaces.CognitiveGameServiceInterface
}

func NewCognitiveHandler(cognitiveService interfaces.CognitiveGameServiceInterface) *CognitiveHandler {
	return &CognitiveHandler{
		BaseHandler:      NewBaseHandler(),
		cognitiveService: cognitiveService,
	}
}

func (h *CognitiveHandler) CreateCognitiveGame(c echo.Context) error {
	var req schemas.CognitiveGameCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	game, err := h.cognitiveService.CreateCognitiveGame(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد بازی شناختی", err.Error())
	}

	return response.Created(c, "بازی شناختی با موفقیت ایجاد شد", game)
}

func (h *CognitiveHandler) GetCognitiveGameByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه بازی الزامی است", "")
	}

	game, err := h.cognitiveService.GetCognitiveGameById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "بازی شناختی")
	}

	return response.OK(c, "بازی شناختی با موفقیت دریافت شد", game)
}

func (h *CognitiveHandler) ListCognitiveGames(c echo.Context) error {
	var req schemas.CognitiveGameListRequest
	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, "بدنه درخواست نامعتبر است", "")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	if h.GetUserRole(c) == "user" {
		req.UserID = h.GetUserID(c)
	}

	games, err := h.cognitiveService.GetAllCognitiveGames(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت بازی‌های شناختی", err)
	}

	return response.OK(c, "بازی‌های شناختی با موفقیت دریافت شدند", games)
}

func (h *CognitiveHandler) UpdateCognitiveGame(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه بازی الزامی است", "")
	}

	var req schemas.CognitiveGameUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	game, err := h.cognitiveService.UpdateCognitiveGame(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی بازی شناختی", err.Error())
	}

	return response.OK(c, "بازی شناختی با موفقیت بروزرسانی شد", game)
}

func (h *CognitiveHandler) DeleteCognitiveGame(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه بازی الزامی است", "")
	}

	if err := h.cognitiveService.DeleteCognitiveGame(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف بازی شناختی", err.Error())
	}

	return response.OK(c, "بازی شناختی با موفقیت حذف شد", nil)
}

func (h *CognitiveHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}
