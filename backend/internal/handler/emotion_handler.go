package handler

import (
	"net/http"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/labstack/echo/v4"
)

type EmotionHandler struct {
	*BaseHandler
	emotionService     interfaces.EmotionTriangleServiceInterface
	stressService      interfaces.StressEventServiceInterface
	bodyTensionService interfaces.BodyTensionServiceInterface
}

func NewEmotionHandler(emotionService interfaces.EmotionTriangleServiceInterface, stressService interfaces.StressEventServiceInterface, bodyTensionService interfaces.BodyTensionServiceInterface) *EmotionHandler {
	return &EmotionHandler{
		BaseHandler:        NewBaseHandler(),
		emotionService:     emotionService,
		stressService:      stressService,
		bodyTensionService: bodyTensionService,
	}
}

func (h *EmotionHandler) CreateEmotionInteraction(c echo.Context) error {
	var req schemas.EmotionInteractionCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	interaction, err := h.emotionService.CreateEmotionInteraction(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد تعامل احساسی", err.Error())
	}

	return response.Created(c, "تعامل احساسی با موفقیت ایجاد شد", interaction)
}

func (h *EmotionHandler) GetEmotionInteractionByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه تعامل الزامی است")
	}

	interaction, err := h.emotionService.GetEmotionInteractionById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "تعامل احساسی")
	}

	return response.OK(c, "تعامل احساسی با موفقیت دریافت شد", interaction)
}

func (h *EmotionHandler) ListEmotionInteractions(c echo.Context) error {
	var req schemas.EmotionInteractionListRequest
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

	interactions, err := h.emotionService.GetAllEmotionInteractions(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت تعاملات احساسی", err)
	}

	return response.OK(c, "تعاملات احساسی با موفقیت دریافت شدند", interactions)
}

func (h *EmotionHandler) UpdateEmotionInteraction(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه تعامل الزامی است")
	}

	var req schemas.EmotionInteractionUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	interaction, err := h.emotionService.UpdateEmotionInteraction(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی تعامل احساسی", err.Error())
	}

	return response.OK(c, "تعامل احساسی با موفقیت بروزرسانی شد", interaction)
}

func (h *EmotionHandler) DeleteEmotionInteraction(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه تعامل الزامی است")
	}

	if err := h.emotionService.DeleteEmotionInteraction(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف تعامل احساسی", err.Error())
	}

	return response.OK(c, "تعامل احساسی با موفقیت حذف شد", nil)
}

func (h *EmotionHandler) CreateStressEvent(c echo.Context) error {
	var req schemas.StressEventCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	event, err := h.stressService.CreateStressEvent(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد رویداد استرس", err.Error())
	}

	return response.Created(c, "رویداد استرس با موفقیت ایجاد شد", event)
}

func (h *EmotionHandler) GetStressEventByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه رویداد الزامی است")
	}

	event, err := h.stressService.GetStressEventById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "رویداد استرس")
	}

	return response.OK(c, "رویداد استرس با موفقیت دریافت شد", event)
}

func (h *EmotionHandler) ListStressEvents(c echo.Context) error {
	var req schemas.StressEventListRequest
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

	events, err := h.stressService.GetAllStressEvents(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت رویدادهای استرس", err)
	}

	return response.OK(c, "رویدادهای استرس با موفقیت دریافت شدند", events)
}

func (h *EmotionHandler) UpdateStressEvent(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه رویداد الزامی است")
	}

	var req schemas.StressEventUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	event, err := h.stressService.UpdateStressEvent(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی رویداد استرس", err.Error())
	}

	return response.OK(c, "رویداد استرس با موفقیت بروزرسانی شد", event)
}

func (h *EmotionHandler) DeleteStressEvent(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه رویداد الزامی است")
	}

	if err := h.stressService.DeleteStressEvent(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف رویداد استرس", err.Error())
	}

	return response.OK(c, "رویداد استرس با موفقیت حذف شد", nil)
}

func (h *EmotionHandler) CreateBodyTensionMap(c echo.Context) error {
	var req schemas.BodyTensionCreateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	req.UserID = h.GetUserID(c)

	tension, err := h.bodyTensionService.CreateBodyTension(c.Request().Context(), &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در ایجاد نقشه تنش بدنی", err.Error())
	}

	return response.Created(c, "نقشه تنش بدنی با موفقیت ایجاد شد", tension)
}

func (h *EmotionHandler) GetBodyTensionMapByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه نقشه الزامی است")
	}

	tension, err := h.bodyTensionService.GetBodyTensionById(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "نقشه تنش بدنی")
	}

	return response.OK(c, "نقشه تنش بدنی با موفقیت دریافت شد", tension)
}

func (h *EmotionHandler) ListBodyTensionMaps(c echo.Context) error {
	var req schemas.BodyTensionListRequest
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

	tensions, err := h.bodyTensionService.GetAllBodyTensions(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت نقشه‌های تنش بدنی", err)
	}

	return response.OK(c, "نقشه‌های تنش بدنی با موفقیت دریافت شدند", tensions)
}

func (h *EmotionHandler) UpdateBodyTensionMap(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه نقشه الزامی است")
	}

	var req schemas.BodyTensionUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	tension, err := h.bodyTensionService.UpdateBodyTension(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی نقشه تنش بدنی", err.Error())
	}

	return response.OK(c, "نقشه تنش بدنی با موفقیت بروزرسانی شد", tension)
}

func (h *EmotionHandler) DeleteBodyTensionMap(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return h.BadRequest(c, "شناسه نقشه الزامی است")
	}

	if err := h.bodyTensionService.DeleteBodyTension(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف نقشه تنش بدنی", err.Error())
	}

	return response.OK(c, "نقشه تنش بدنی با موفقیت حذف شد", nil)
}

func (h *EmotionHandler) BadRequest(c echo.Context, message string) error {
	return response.BadRequest(c, message, "")
}

func (h *EmotionHandler) NotFound(c echo.Context, resource string) error {
	return response.Error(c, http.StatusNotFound, resource+" یافت نشد", "")
}

func (h *EmotionHandler) InternalError(c echo.Context, message string, err error) error {
	return response.InternalServerError(c, message, err.Error())
}

func (h *EmotionHandler) GetUserRole(c echo.Context) string {
	role, ok := c.Get("user_role").(string)
	if !ok {
		return ""
	}
	return role
}
