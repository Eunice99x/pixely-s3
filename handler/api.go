package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/eunice99x/pixely/kyc/service"
	"github.com/eunice99x/pixely/kyc/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	db    *sqlx.DB
	store *storage.S3Bucket
}

func New(db *sqlx.DB, store *storage.S3Bucket) *Handler {
	return &Handler{db: db, store: store}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "missing image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	meta, err := service.ExtractMetadata(data)
	if err != nil {
		http.Error(w, "not a valid image", http.StatusBadRequest)
		return
	}

	meta.Key = uuid.NewString() + "." + meta.Format
	if err := h.store.Upload(r.Context(), meta.Key, data); err != nil {
		http.Error(w, "upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.db.QueryRow(
		`INSERT INTO images (key, format, width, height, size, make, model, taken_at, lat, lng)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		meta.Key, meta.Format, meta.Width, meta.Height, meta.Size,
		meta.Make, meta.Model, meta.TakenAt, meta.Lat, meta.Lng,
	).Scan(&meta.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metaJSON, _ := json.Marshal(meta)
	if err := h.store.Upload(r.Context(), meta.Key+".json", metaJSON); err != nil {
		http.Error(w, "metadata upload failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(meta)
}

func (h *Handler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	var meta service.Metadata
	err := h.db.Get(&meta, `SELECT * FROM images WHERE id = $1`, chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(meta)
}
