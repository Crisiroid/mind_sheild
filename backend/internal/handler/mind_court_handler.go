package handler

import (
	"net/http"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type MindCourtHandler struct {
	*BaseHandler
	mindCourtService interfaces.MindCourtServiceInterface
}

func NewMindCourtHandler(mindCourtService interfaces.MindCourtServiceInterface) *MindCourtHandler {
	return &MindCourtHandler{
		BaseHandler:      NewBaseHandler(),
		mindCourtService: mindCourtService,
	}
}

func (h *MindCourtHandler) CreateMindCourtEvidence(c echo.Context) error {
	var req schemas.MindCourtCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	evidence, err := h.mindCourtService.CreateMindCourt(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد شواهد دادگاه ذهنی", err.Error())
	}

	return response.Created(c, "شواهد دادگاه ذهنی با موفقیت ایجاد شد", evidence)
}

func (h *MindCourtHandler) GetMindCourtEvidenceByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه شواهد الزامی است", "")
	}

	evidence, err := h.mindCourtService.GetMindCourtById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "شواهد دادگاه ذهنی")
	}

	return response.OK(c, "شواهد دادگاه ذهنی با موفقیت دریافت شد", evidence)
}

func (h *MindCourtHandler) ListMindCourtEvidence(c echo.Context) error {
	var req schemas.MindCourtListRequest
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

	evidence, err := h.mindCourtService.GetAllMindCourts(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت شواهد دادگاه ذهنی", err)
	}

	return response.OK(c, "شواهد دادگاه ذهنی با موفقیت دریافت شدند", evidence)
}

func (h *MindCourtHandler) UpdateMindCourtEvidence(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه شواهد الزامی است", "")
	}

	var req schemas.MindCourtUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	evidence, err := h.mindCourtService.UpdateMindCourt(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی شواهد دادگاه ذهنی", err.Error())
	}

	return response.OK(c, "شواهد دادگاه ذهنی با موفقیت بروزرسانی شد", evidence)
}

func (h *MindCourtHandler) DeleteMindCourtEvidence(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه شواهد الزامی است", "")
	}

	if err := h.mindCourtService.DeleteMindCourt(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف شواهد دادگاه ذهنی", err.Error())
	}

	return response.OK(c, "شواهد دادگاه ذهنی با موفقیت حذف شد", nil)
}

func (h *MindCourtHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}
