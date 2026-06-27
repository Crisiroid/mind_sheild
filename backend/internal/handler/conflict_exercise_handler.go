package handler

import (
	"net/http"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type ConflictExerciseHandler struct {
	*BaseHandler
	conflictService interfaces.ConflictExerciseServiceInterface
}

func NewConflictExerciseHandler(conflictService interfaces.ConflictExerciseServiceInterface) *ConflictExerciseHandler {
	return &ConflictExerciseHandler{
		BaseHandler:     NewBaseHandler(),
		conflictService: conflictService,
	}
}

func (h *ConflictExerciseHandler) CreateConflictExercise(c echo.Context) error {
	var req schemas.ConflictExerciseCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	exercise, err := h.conflictService.CreateConflictExercise(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد تمرین تعارض", err.Error())
	}

	return response.Created(c, "تمرین تعارض با موفقیت ایجاد شد", exercise)
}

func (h *ConflictExerciseHandler) GetConflictExerciseByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه تمرین الزامی است", "")
	}

	exercise, err := h.conflictService.GetConflictExerciseById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "تمرین تعارض")
	}

	return response.OK(c, "تمرین تعارض با موفقیت دریافت شد", exercise)
}

func (h *ConflictExerciseHandler) ListConflictExercises(c echo.Context) error {
	var req schemas.ConflictExerciseListRequest
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

	exercises, err := h.conflictService.GetAllConflictExercises(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت تمرین‌های تعارض", err)
	}

	return response.OK(c, "تمرین‌های تعارض با موفقیت دریافت شدند", exercises)
}

func (h *ConflictExerciseHandler) UpdateConflictExercise(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه تمرین الزامی است", "")
	}

	var req schemas.ConflictExerciseUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	exercise, err := h.conflictService.UpdateConflictExercise(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی تمرین تعارض", err.Error())
	}

	return response.OK(c, "تمرین تعارض با موفقیت بروزرسانی شد", exercise)
}

func (h *ConflictExerciseHandler) DeleteConflictExercise(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه تمرین الزامی است", "")
	}

	if err := h.conflictService.DeleteConflictExercise(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف تمرین تعارض", err.Error())
	}

	return response.OK(c, "تمرین تعارض با موفقیت حذف شد", nil)
}

func (h *ConflictExerciseHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}
