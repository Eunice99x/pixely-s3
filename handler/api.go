package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/eunice99x/pixely/kyc/service"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Handler {
	return &Handler{db: db}
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

	err = h.db.QueryRow(
		`INSERT INTO images (format, width, height, size) VALUES ($1, $2, $3, $4) RETURNING id`,
		meta.Format, meta.Width, meta.Height, meta.Size,
	).Scan(&meta.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
