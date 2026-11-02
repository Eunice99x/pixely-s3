package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/eunice99x/pixely/kyc/internal/database"
	"github.com/eunice99x/pixely/kyc/internal/model"
	"github.com/eunice99x/pixely/kyc/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type ImageStore interface {
	CreateImage(ctx context.Context, img *model.Image) error
	GetImage(ctx context.Context, id uint) (model.Image, error)
	Ping(ctx context.Context) error
}

type FileStore interface {
	Upload(ctx context.Context, key string, data []byte) error
}

type Handler struct {
	db    ImageStore
	files FileStore
}

func New(db ImageStore, files FileStore) *Handler {
	return &Handler{db: db, files: files}
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
	if err := h.db.CreateImage(ctx, &img); err != nil {
		return err
	}

	metaJSON, _ := json.Marshal(img)
	if err := h.files.Upload(ctx, img.Key, metaJSON); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "metadata upload failed: "+err.Error())
	}

	return c.JSON(http.StatusOK, img)
}

func (h *Handler) GetMetadata(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}

	img, err := h.db.GetImage(c.Request().Context(), uint(id))
	if errors.Is(err, database.ErrNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, img)
}

func (h *Handler) Health(c *echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		return c.JSON(http.StatusServiceUnavailable, model.HealthResponse{Status: "unavailable", Database: "down"})
	}
	return c.JSON(http.StatusOK, model.HealthResponse{Status: "ok", Database: "up"})
}
