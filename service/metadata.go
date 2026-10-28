package service

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"time"

	"github.com/evanoberholster/imagemeta"
)

type Metadata struct {
	ID      int        `db:"id" json:"id"`
	Key     string     `db:"key" json:"key"`
	Format  string     `db:"format" json:"format"`
	Width   int        `db:"width" json:"width"`
	Height  int        `db:"height" json:"height"`
	Size    int        `db:"size" json:"size"`
	Make    string     `db:"make" json:"make"`
	Model   string     `db:"model" json:"model"`
	TakenAt *time.Time `db:"taken_at" json:"taken_at"`
	Lat     *float64   `db:"lat" json:"lat"`
	Lng     *float64   `db:"lng" json:"lng"`
}

func ExtractMetadata(data []byte) (Metadata, error) {
	m := Metadata{Size: len(data)}

	ex, exErr := imagemeta.Decode(bytes.NewReader(data))
	if exErr == nil {
		m.Format = ex.ImageType.String()
		m.Width, m.Height = int(ex.ExifIFD.PixelXDimension), int(ex.ExifIFD.PixelYDimension)

		m.Make, m.Model = ex.IFD0.Make, ex.IFD0.Model

		if t := ex.ExifIFD.DateTimeOriginal; !t.IsZero() {
			m.TakenAt = &t
		}
		if lat, lng := ex.GPS.Latitude(), ex.GPS.Longitude(); lat != 0 || lng != 0 {
			m.Lat, m.Lng = &lat, &lng
		}
	}

	if cfg, format, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
		m.Format, m.Width, m.Height = format, cfg.Width, cfg.Height
	} else if exErr != nil {
		return Metadata{}, err
	}

	return m, nil
}
