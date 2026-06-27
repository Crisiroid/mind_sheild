package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"psychology-backend/internal/interfaces"
	"psychology-backend/pkg/response"
	"psychology-backend/pkg/schemas"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type WeeklyMediaContentHandler struct {
	*BaseHandler
	mediaService    interfaces.WeeklyMediaContentServiceInterface
	uploadDirectory string
	baseURL         string
}

func NewWeeklyMediaContentHandler(mediaService interfaces.WeeklyMediaContentServiceInterface, uploadDirectory, baseURL string) *WeeklyMediaContentHandler {
	return &WeeklyMediaContentHandler{
		BaseHandler:     NewBaseHandler(),
		mediaService:    mediaService,
		uploadDirectory: uploadDirectory,
		baseURL:         baseURL,
	}
}

func (h *WeeklyMediaContentHandler) UploadMediaContent(c echo.Context) error {
	adminID := h.GetUserID(c)
	if adminID == "" {
		return response.BadRequest(c, "شناسه ادمین الزامی است", "")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return response.BadRequest(c, "فرم درخواست نامعتبر است", err.Error())
	}

	files := form.File["file"]
	if len(files) == 0 {
		return response.BadRequest(c, "فایل الزامی است", "")
	}

	file := files[0]
	if file.Size > 100*1024*1024 {
		return response.BadRequest(c, "حجم فایل نباید بیشتر از ۱۰۰ مگابایت باشد", "")
	}

	weekNumberStr := form.Value["week_number"]
	if len(weekNumberStr) == 0 {
		return response.BadRequest(c, "شماره هفته الزامی است", "")
	}

	weekNumber, err := strconv.Atoi(weekNumberStr[0])
	if err != nil || weekNumber < 1 || weekNumber > 52 {
		return response.BadRequest(c, "شماره هفته باید بین ۱ تا ۵۲ باشد", "")
	}

	fileType := ""
	if len(form.Value["file_type"]) > 0 {
		fileType = form.Value["file_type"][0]
	} else {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		fileType = detectFileType(ext)
	}

	description := ""
	if len(form.Value["description"]) > 0 {
		description = form.Value["description"][0]
	}

	src, err := file.Open()
	if err != nil {
		return h.InternalError(c, "خطا در باز کردن فایل", err)
	}
	defer src.Close()

	uploadPath := filepath.Join(h.uploadDirectory, "weekly-media")
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return h.InternalError(c, "خطا در ایجاد پوشه آپلود", err)
	}

	ext := filepath.Ext(file.Filename)
	nameWithoutExt := strings.TrimSuffix(file.Filename, ext)
	uniqueFileName := fmt.Sprintf("%s_%s%s", nameWithoutExt, uuid.New().String()[:8], ext)

	destPath := filepath.Join(uploadPath, uniqueFileName)
	dst, err := os.Create(destPath)
	if err != nil {
		return h.InternalError(c, "خطا در ذخیره فایل", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return h.InternalError(c, "خطا در ذخیره فایل", err)
	}

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = detectContentType(ext)
	}

	req := &schemas.WeeklyMediaContentCreateRequest{
		WeekNumber:  weekNumber,
		FileType:    fileType,
		Description: description,
	}

	storagePath := fmt.Sprintf("weekly-media/%s", uniqueFileName)
	media, err := h.mediaService.UploadMediaContent(
		c.Request().Context(),
		file.Filename,
		fileType,
		contentType,
		storagePath,
		file.Size,
		req,
		adminID,
	)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در آپلود فایل", err.Error())
	}

	return response.Created(c, "فایل با موفقیت آپلود شد", media)
}

func (h *WeeklyMediaContentHandler) GetMediaContentByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه محتوا الزامی است", "")
	}

	media, err := h.mediaService.GetMediaContentByID(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "محتوای مدیا")
	}

	return response.OK(c, "محتوای مدیا با موفقیت دریافت شد", media)
}

func (h *WeeklyMediaContentHandler) ListMediaContent(c echo.Context) error {
	var req schemas.WeeklyMediaContentListRequest
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
		isActive := true
		req.IsActive = &isActive
	}

	mediaList, err := h.mediaService.GetAllMediaContent(c.Request().Context(), &req)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت محتوای مدیا", err)
	}

	return response.OK(c, "محتوای مدیا با موفقیت دریافت شد", mediaList)
}

func (h *WeeklyMediaContentHandler) UpdateMediaContent(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه محتوا الزامی است", "")
	}

	var req schemas.WeeklyMediaContentUpdateRequest
	if err := h.BindAndValidate(c, &req); err != nil {
		return err
	}

	media, err := h.mediaService.UpdateMediaContent(c.Request().Context(), id, &req)
	if err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در بروزرسانی محتوا", err.Error())
	}

	return response.OK(c, "محتوای مدیا با موفقیت بروزرسانی شد", media)
}

func (h *WeeklyMediaContentHandler) DeleteMediaContent(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه محتوا الزامی است", "")
	}

	if err := h.mediaService.DeleteMediaContent(c.Request().Context(), id); err != nil {
		return response.Error(c, http.StatusBadRequest, "خطا در حذف محتوا", err.Error())
	}

	return response.Success(c, http.StatusOK, "محتوای مدیا با موفقیت حذف شد", nil)
}

func (h *WeeklyMediaContentHandler) GetMediaContentByWeek(c echo.Context) error {
	weekNumberStr := c.Param("week_number")
	weekNumber, err := strconv.Atoi(weekNumberStr)
	if err != nil || weekNumber < 1 || weekNumber > 52 {
		return response.BadRequest(c, "شماره هفته باید بین ۱ تا ۵۲ باشد", "")
	}

	pageStr := c.QueryParam("page")
	pageSizeStr := c.QueryParam("page_size")

	page := 1
	pageSize := 20

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	mediaList, err := h.mediaService.GetMediaContentByWeek(c.Request().Context(), weekNumber, page, pageSize)
	if err != nil {
		return h.InternalError(c, "خطا در دریافت محتوای مدیا", err)
	}

	return response.OK(c, "محتوای مدیا با موفقیت دریافت شد", mediaList)
}

func (h *WeeklyMediaContentHandler) DownloadMediaContent(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return response.BadRequest(c, "شناسه محتوا الزامی است", "")
	}

	media, err := h.mediaService.GetMediaContentByID(c.Request().Context(), id)
	if err != nil {
		return h.NotFound(c, "محتوای مدیا")
	}

	if !media.IsActive {
		return response.Error(c, http.StatusForbidden, "این محتوا در دسترس نیست", "")
	}

	if err := h.mediaService.IncrementDownloadCount(c.Request().Context(), id); err != nil {
		fmt.Printf("Error incrementing download count: %v\n", err)
	}

	filePath := filepath.Join(h.uploadDirectory, media.StoragePath)
	return c.Attachment(filePath, media.OriginalName)
}

func detectFileType(ext string) string {
	audioExts := map[string]bool{".mp3": true, ".wav": true, ".ogg": true, ".m4a": true, ".aac": true}
	videoExts := map[string]bool{".mp4": true, ".avi": true, ".mkv": true, ".mov": true, ".webm": true}
	imageExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".svg": true}
	docExts := map[string]bool{".pdf": true, ".doc": true, ".docx": true, ".txt": true, ".xlsx": true, ".pptx": true}

	if audioExts[ext] {
		return "audio"
	}
	if videoExts[ext] {
		return "video"
	}
	if imageExts[ext] {
		return "image"
	}
	if docExts[ext] {
		return "document"
	}
	return "document"
}

func detectContentType(ext string) string {
	contentTypes := map[string]string{
		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".ogg":  "audio/ogg",
		".mp4":  "video/mp4",
		".avi":  "video/x-msvideo",
		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
	}

	if ct, ok := contentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}
