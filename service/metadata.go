package service

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
)

type Metadata struct {
	ID     int    `db:"id" json:"id"`
	Format string `db:"format" json:"format"`
	Width  int    `db:"width" json:"width"`
	Height int    `db:"height" json:"height"`
	Size   int    `db:"size" json:"size"`
}

func ExtractMetadata(data []byte) (Metadata, error) {
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return Metadata{}, err
	}
	return Metadata{Format: format, Width: cfg.Width, Height: cfg.Height, Size: len(data)}, nil
}
