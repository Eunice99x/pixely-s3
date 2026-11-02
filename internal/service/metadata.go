package service

import (
	"bytes"
	"encoding/json"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/eunice99x/pixely/kyc/internal/model"
	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/meta/exif"
)

func ExtractMetadata(data []byte) (model.Image, error) {
	m := model.Image{Size: len(data)}

	ex, exErr := imagemeta.Decode(bytes.NewReader(data))
	if exErr == nil {
		fillFromExif(&m, ex)
	}

	if cfg, format, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
		m.Format, m.Width, m.Height = format, cfg.Width, cfg.Height
	} else if exErr != nil {
		return model.Image{}, err
	}

	return m, nil
}

func fillFromExif(m *model.Image, ex exif.Exif) {
	ifd0, e := ex.IFD0, ex.ExifIFD

	m.Format = ex.ImageType.Extension()
	m.MimeType = ex.ImageType.MIMEType().String()
	m.Width, m.Height = int(e.PixelXDimension), int(e.PixelYDimension)
	if m.Width == 0 || m.Height == 0 {
		m.Width, m.Height = int(ifd0.ImageWidth), int(ifd0.ImageHeight)
	}

	m.Make, m.Model, m.CameraSerial = ifd0.Make, ifd0.Model, ex.CameraSerial
	if m.CameraSerial == "" {
		m.CameraSerial = e.BodySerialNumber
	}
	m.LensMake, m.LensModel = e.LensMake, e.LensModel
	m.Software, m.Artist, m.Copyright = ifd0.Software, ifd0.Artist, ifd0.Copyright
	if ifd0.Orientation != 0 {
		m.Orientation = ifd0.Orientation.String()
	}

	if t := e.DateTimeOriginal; !t.IsZero() {
		m.TakenAt = &t
	} else if t := e.CreateDate; !t.IsZero() {
		m.TakenAt = &t
	}
	if t := ifd0.ModifyDate; !t.IsZero() {
		m.ModifiedAt = &t
	}

	if e.ExposureTime != 0 {
		m.ExposureTime = e.ExposureTime.String()
	}
	m.FNumber = float32(e.FNumber)
	m.ISO = e.ISOSpeedRatings
	m.FocalLength = float32(e.FocalLength)
	m.FocalLength35mm = float32(e.FocalLengthIn35mmFormat)
	if e.ExposureBias != 0 {
		m.ExposureBias = e.ExposureBias.String()
	}
	m.FlashFired = e.Flash.Fired()

	if lat, lng := ex.GPS.Latitude(), ex.GPS.Longitude(); lat != 0 || lng != 0 {
		alt := ex.GPS.Altitude()
		m.Lat, m.Lng, m.Altitude = &lat, &lng, &alt
	}

	if raw, err := json.Marshal(ex); err == nil {
		m.Exif = raw
	}
}
