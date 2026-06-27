package handler

import (
	"net/http"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type MoodTrackerHandler struct {
	*BaseHandler
	moodTrackerService interfaces.MoodTrackerServiceInterface
}

func NewMoodTrackerHandler(moodTrackerService interfaces.MoodTrackerServiceInterface) *MoodTrackerHandler {
	return &MoodTrackerHandler{
		BaseHandler:        NewBaseHandler(),
		moodTrackerService: moodTrackerService,
	}
}

func (h *MoodTrackerHandler) CreateMoodTracker(c echo.Context) error {
	var req schemas.MoodTrackerCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	mood, err := h.moodTrackerService.CreateMoodTracker(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد ردیاب خلق", err.Error())
	}

	return response.Created(c, "ردیاب خلق با موفقیت ایجاد شد", mood)
}

func (h *MoodTrackerHandler) GetMoodTrackerByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه ردیاب الزامی است", "")
	}

	mood, err := h.moodTrackerService.GetMoodTrackerById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "ردیاب خلق")
	}

	return response.OK(c, "ردیاب خلق با موفقیت دریافت شد", mood)
}

func (h *MoodTrackerHandler) ListMoodTrackers(c echo.Context) error {
	var req schemas.MoodTrackerListRequest
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

	moods, err := h.moodTrackerService.GetAllMoodTrackers(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت ردیاب‌های خلق", err)
	}

	return response.OK(c, "ردیاب‌های خلق با موفقیت دریافت شدند", moods)
}

func (h *MoodTrackerHandler) UpdateMoodTracker(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه ردیاب الزامی است", "")
	}

	var req schemas.MoodTrackerUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	mood, err := h.moodTrackerService.UpdateMoodTracker(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی ردیاب خلق", err.Error())
	}

	return response.OK(c, "ردیاب خلق با موفقیت بروزرسانی شد", mood)
}

func (h *MoodTrackerHandler) DeleteMoodTracker(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه ردیاب الزامی است", "")
	}

	if err := h.moodTrackerService.DeleteMoodTracker(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف ردیاب خلق", err.Error())
	}

	return response.OK(c, "ردیاب خلق با موفقیت حذف شد", nil)
}

func (h *MoodTrackerHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}
