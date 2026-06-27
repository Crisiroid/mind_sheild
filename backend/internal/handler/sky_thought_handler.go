package handler

import (
	"net/http"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type SkyThoughtHandler struct {
	*BaseHandler
	skyThoughtService interfaces.SkyThoughtServiceInterface
}

func NewSkyThoughtHandler(skyThoughtService interfaces.SkyThoughtServiceInterface) *SkyThoughtHandler {
	return &SkyThoughtHandler{
		BaseHandler:       NewBaseHandler(),
		skyThoughtService: skyThoughtService,
	}
}

func (h *SkyThoughtHandler) CreateSkyThought(c echo.Context) error {
	var req schemas.SkyThoughtCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	skyThought, err := h.skyThoughtService.CreateSkyThought(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد فکر آسمانی", err.Error())
	}

	return response.Created(c, "فکر آسمانی با موفقیت ایجاد شد", skyThought)
}

func (h *SkyThoughtHandler) GetSkyThoughtByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه فکر الزامی است", "")
	}

	skyThought, err := h.skyThoughtService.GetSkyThoughtById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "فکر آسمانی")
	}

	return response.OK(c, "فکر آسمانی با موفقیت دریافت شد", skyThought)
}

func (h *SkyThoughtHandler) ListSkyThoughts(c echo.Context) error {
	var req schemas.SkyThoughtListRequest
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

	skyThoughts, err := h.skyThoughtService.GetAllSkyThoughts(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت افکار آسمانی", err)
	}

	return response.OK(c, "افکار آسمانی با موفقیت دریافت شدند", skyThoughts)
}

func (h *SkyThoughtHandler) UpdateSkyThought(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه فکر الزامی است", "")
	}

	var req schemas.SkyThoughtUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	skyThought, err := h.skyThoughtService.UpdateSkyThought(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی فکر آسمانی", err.Error())
	}

	return response.OK(c, "فکر آسمانی با موفقیت بروزرسانی شد", skyThought)
}

func (h *SkyThoughtHandler) DeleteSkyThought(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه فکر الزامی است", "")
	}

	if err := h.skyThoughtService.DeleteSkyThought(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف فکر آسمانی", err.Error())
	}

	return response.OK(c, "فکر آسمانی با موفقیت حذف شد", nil)
}

func (h *SkyThoughtHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}
