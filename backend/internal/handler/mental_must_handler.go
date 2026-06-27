package handler

import (
	"net/http"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type MentalMustHandler struct {
	*BaseHandler
	mentalMustService interfaces.MentalMustServiceInterface
}

func NewMentalMustHandler(mentalMustService interfaces.MentalMustServiceInterface) *MentalMustHandler {
	return &MentalMustHandler{
		BaseHandler:       NewBaseHandler(),
		mentalMustService: mentalMustService,
	}
}

func (h *MentalMustHandler) CreateMentalMust(c echo.Context) error {
	var req schemas.MentalMustCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	mentalMust, err := h.mentalMustService.CreateMentalMust(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد باید ذهنی", err.Error())
	}

	return response.Created(c, "باید ذهنی با موفقیت ایجاد شد", mentalMust)
}

func (h *MentalMustHandler) GetMentalMustByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه باید ذهنی الزامی است", "")
	}

	mentalMust, err := h.mentalMustService.GetMentalMustById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "باید ذهنی")
	}

	return response.OK(c, "باید ذهنی با موفقیت دریافت شد", mentalMust)
}

func (h *MentalMustHandler) ListMentalMusts(c echo.Context) error {
	var req schemas.MentalMustListRequest
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

	mentalMusts, err := h.mentalMustService.GetAllMentalMusts(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت بایدهای ذهنی", err)
	}

	return response.OK(c, "بایدهای ذهنی با موفقیت دریافت شدند", mentalMusts)
}

func (h *MentalMustHandler) UpdateMentalMust(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه باید ذهنی الزامی است", "")
	}

	var req schemas.MentalMustUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	mentalMust, err := h.mentalMustService.UpdateMentalMust(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی باید ذهنی", err.Error())
	}

	return response.OK(c, "باید ذهنی با موفقیت بروزرسانی شد", mentalMust)
}

func (h *MentalMustHandler) DeleteMentalMust(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه باید ذهنی الزامی است", "")
	}

	if err := h.mentalMustService.DeleteMentalMust(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف باید ذهنی", err.Error())
	}

	return response.OK(c, "باید ذهنی با موفقیت حذف شد", nil)
}

func (h *MentalMustHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}
