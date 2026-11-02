package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/eunice99x/pixely/kyc/service"
	"github.com/eunice99x/pixely/kyc/storage"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type Handler struct {
	db    *gorm.DB
	store *storage.S3Bucket
}

func New(db *gorm.DB, store *storage.S3Bucket) *Handler {
	return &Handler{db: db, store: store}
}

func (h *Handler) Upload(c *echo.Context) error {
	fh, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing image")
	}
	file, err := fh.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	img, err := service.ExtractMetadata(data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "not a valid image")
	}

	ctx := c.Request().Context()
	img.Key = uuid.NewString() + ".json"
	if err := h.db.WithContext(ctx).Create(&img).Error; err != nil {
		return err
	}

	metaJSON, _ := json.Marshal(img)
	if err := h.store.Upload(ctx, img.Key, metaJSON); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "metadata upload failed: "+err.Error())
	}

	return c.JSON(http.StatusOK, img)
}

func (h *Handler) GetMetadata(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	var img service.Image
	err = h.db.WithContext(c.Request().Context()).First(&img, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, img)
}
