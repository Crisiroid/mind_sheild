package handler

import (
	"net/http"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type BreathingHandler struct {
	*BaseHandler
	breathingService interfaces.BreathingServiceInterface
}

func NewBreathingHandler(breathingService interfaces.BreathingServiceInterface) *BreathingHandler {
	return &BreathingHandler{
		BaseHandler:      NewBaseHandler(),
		breathingService: breathingService,
	}
}

func (h *BreathingHandler) CreateBreathingSession(c echo.Context) error {
	var req schemas.BreathingSessionCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	session, err := h.breathingService.CreateBreathingSession(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد جلسه تنفس", err.Error())
	}

	return response.Created(c, "جلسه تنفس با موفقیت ایجاد شد", session)
}

func (h *BreathingHandler) GetBreathingSessionByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه جلسه الزامی است", "")
	}

	session, err := h.breathingService.GetBreathingSessionById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "جلسه تنفس")
	}

	return response.OK(c, "جلسه تنفس با موفقیت دریافت شد", session)
}

func (h *BreathingHandler) ListBreathingSessions(c echo.Context) error {
	var req schemas.BreathingSessionListRequest
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

	sessions, err := h.breathingService.GetAllBreathingSessions(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت جلسات تنفس", err)
	}

	return response.OK(c, "جلسات تنفس با موفقیت دریافت شدند", sessions)
}

func (h *BreathingHandler) UpdateBreathingSession(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه جلسه الزامی است", "")
	}

	var req schemas.BreathingSessionUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	session, err := h.breathingService.UpdateBreathingSession(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی جلسه تنفس", err.Error())
	}

	return response.OK(c, "جلسه تنفس با موفقیت بروزرسانی شد", session)
}

func (h *BreathingHandler) DeleteBreathingSession(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه جلسه الزامی است", "")
	}

	if err := h.breathingService.DeleteBreathingSession(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف جلسه تنفس", err.Error())
	}

	return response.OK(c, "جلسه تنفس با موفقیت حذف شد", nil)
}

func (h *BreathingHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}
