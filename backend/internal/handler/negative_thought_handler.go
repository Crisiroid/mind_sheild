package handler

import (
	"net/http"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type NegativeThoughtHandler struct {
	*BaseHandler
	negativeThoughtService interfaces.NegativeThoughtServiceInterface
}

func NewNegativeThoughtHandler(negativeThoughtService interfaces.NegativeThoughtServiceInterface) *NegativeThoughtHandler {
	return &NegativeThoughtHandler{
		BaseHandler:            NewBaseHandler(),
		negativeThoughtService: negativeThoughtService,
	}
}

func (h *NegativeThoughtHandler) CreateNegativeThought(c echo.Context) error {
	var req schemas.NegativeThoughtCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	thought, err := h.negativeThoughtService.CreateNegativeThought(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد فکر منفی", err.Error())
	}

	return response.Created(c, "فکر منفی با موفقیت ایجاد شد", thought)
}

func (h *NegativeThoughtHandler) GetNegativeThoughtByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه فکر الزامی است", "")
	}

	thought, err := h.negativeThoughtService.GetNegativeThoughtById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "فکر منفی")
	}

	return response.OK(c, "فکر منفی با موفقیت دریافت شد", thought)
}

func (h *NegativeThoughtHandler) ListNegativeThoughts(c echo.Context) error {
	var req schemas.NegativeThoughtListRequest
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

	thoughts, err := h.negativeThoughtService.GetAllNegativeThoughts(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت افکار منفی", err)
	}

	return response.OK(c, "افکار منفی با موفقیت دریافت شدند", thoughts)
}

func (h *NegativeThoughtHandler) UpdateNegativeThought(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه فکر الزامی است", "")
	}

	var req schemas.NegativeThoughtUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	thought, err := h.negativeThoughtService.UpdateNegativeThought(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی فکر منفی", err.Error())
	}

	return response.OK(c, "فکر منفی با موفقیت بروزرسانی شد", thought)
}

func (h *NegativeThoughtHandler) DeleteNegativeThought(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه فکر الزامی است", "")
	}

	if err := h.negativeThoughtService.DeleteNegativeThought(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف فکر منفی", err.Error())
	}

	return response.OK(c, "فکر منفی با موفقیت حذف شد", nil)
}

func (h *NegativeThoughtHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}
