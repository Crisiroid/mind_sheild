package handler

import (
	"net/http"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type MindfulnessHandler struct {
	*BaseHandler
	mindfulTimerService interfaces.MindfulTimerServiceInterface
	acceptanceService   interfaces.AcceptanceServiceInterface
}

func NewMindfulnessHandler(mindfulTimerService interfaces.MindfulTimerServiceInterface, acceptanceService interfaces.AcceptanceServiceInterface) *MindfulnessHandler {
	return &MindfulnessHandler{
		BaseHandler:         NewBaseHandler(),
		mindfulTimerService: mindfulTimerService,
		acceptanceService:   acceptanceService,
	}
}

func (h *MindfulnessHandler) CreateMindfulTimer(c echo.Context) error {
	var req schemas.MindfulTimerCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	timer, err := h.mindfulTimerService.CreateMindfulTimer(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد تایمر ذهن‌آگاهی", err.Error())
	}

	return response.Created(c, "تایمر ذهن‌آگاهی با موفقیت ایجاد شد", timer)
}

func (h *MindfulnessHandler) GetMindfulTimerByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه تایمر الزامی است")
	}

	timer, err := h.mindfulTimerService.GetMindfulTimerById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "تایمر ذهن‌آگاهی")
	}

	return response.OK(c, "تایمر ذهن‌آگاهی با موفقیت دریافت شد", timer)
}

func (h *MindfulnessHandler) ListMindfulTimers(c echo.Context) error {
	var req schemas.MindfulTimerListRequest
	if err := c.Bind(&req); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
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

	timers, err := h.mindfulTimerService.GetAllMindfulTimers(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت تایمرهای ذهن‌آگاهی", err)
	}

	return response.OK(c, "تایمرهای ذهن‌آگاهی با موفقیت دریافت شدند", timers)
}

func (h *MindfulnessHandler) UpdateMindfulTimer(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه تایمر الزامی است")
	}

	var req schemas.MindfulTimerUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	timer, err := h.mindfulTimerService.UpdateMindfulTimer(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی تایمر ذهن‌آگاهی", err.Error())
	}

	return response.OK(c, "تایمر ذهن‌آگاهی با موفقیت بروزرسانی شد", timer)
}

func (h *MindfulnessHandler) DeleteMindfulTimer(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه تایمر الزامی است")
	}

	if err := h.mindfulTimerService.DeleteMindfulTimer(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف تایمر ذهن‌آگاهی", err.Error())
	}

	return response.OK(c, "تایمر ذهن‌آگاهی با موفقیت حذف شد", nil)
}

func (h *MindfulnessHandler) CreateAcceptanceExercise(c echo.Context) error {
	var req schemas.AcceptanceExerciseCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	exercise, err := h.acceptanceService.CreateAcceptanceExercise(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد تمرین پذیرش", err.Error())
	}

	return response.Created(c, "تمرین پذیرش با موفقیت ایجاد شد", exercise)
}

func (h *MindfulnessHandler) GetAcceptanceExerciseByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه تمرین الزامی است")
	}

	exercise, err := h.acceptanceService.GetAcceptanceExerciseById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "تمرین پذیرش")
	}

	return response.OK(c, "تمرین پذیرش با موفقیت دریافت شد", exercise)
}

func (h *MindfulnessHandler) ListAcceptanceExercises(c echo.Context) error {
	var req schemas.AcceptanceExerciseListRequest
	if err := c.Bind(&req); err != nil {
		return h.BadRequest(c, "بدنه درخواست نامعتبر است")
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

	exercises, err := h.acceptanceService.GetAllAcceptanceExercises(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت تمرین‌های پذیرش", err)
	}

	return response.OK(c, "تمرین‌های پذیرش با موفقیت دریافت شدند", exercises)
}

func (h *MindfulnessHandler) UpdateAcceptanceExercise(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه تمرین الزامی است")
	}

	var req schemas.AcceptanceExerciseUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	exercise, err := h.acceptanceService.UpdateAcceptanceExercise(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی تمرین پذیرش", err.Error())
	}

	return response.OK(c, "تمرین پذیرش با موفقیت بروزرسانی شد", exercise)
}

func (h *MindfulnessHandler) DeleteAcceptanceExercise(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه تمرین الزامی است")
	}

	if err := h.acceptanceService.DeleteAcceptanceExercise(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف تمرین پذیرش", err.Error())
	}

	return response.OK(c, "تمرین پذیرش با موفقیت حذف شد", nil)
}

func (h *MindfulnessHandler) BadRequest(c echo.Context, message string) error {
	return response.BadRequest(c, message, "")
}

func (h *MindfulnessHandler) NotFound(c echo.Context, resource string) error {
	return response.Error(c, http.StatusNotFound, resource+" یافت نشد", "")
}

func (h *MindfulnessHandler) InternalError(c echo.Context, message string, err error) error {
	return response.InternalServerError(c, message, err.Error())
}

func (h *MindfulnessHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}
