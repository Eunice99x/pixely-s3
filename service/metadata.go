package service

import (
	"bytes"
	"encoding/json"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"time"

	"github.com/evanoberholster/imagemeta"
	"github.com/evanoberholster/imagemeta/meta/exif"
	"gorm.io/datatypes"
)

// Image holds the most useful EXIF fields as columns (so they can be queried)
// and the full decoded EXIF, including vendor maker notes, in the Exif column.
type Image struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Key      string `gorm:"uniqueIndex;not null" json:"key"`
	Format   string `json:"format"`
	MimeType string `json:"mime_type"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`

	// camera / author
	Make         string `json:"make"`
	Model        string `json:"model"`
	CameraSerial string `json:"camera_serial"`
	LensMake     string `json:"lens_make"`
	LensModel    string `json:"lens_model"`
	Software     string `json:"software"`
	Artist       string `json:"artist"`
	Copyright    string `json:"copyright"`
	Orientation  string `json:"orientation"`

	// time
	TakenAt    *time.Time `gorm:"index" json:"taken_at"`
	ModifiedAt *time.Time `json:"modified_at"`

	// exposure
	ExposureTime    string  `json:"exposure_time"`
	FNumber         float32 `json:"f_number"`
	ISO             uint32  `json:"iso"`
	FocalLength     float32 `json:"focal_length"`
	FocalLength35mm float32 `json:"focal_length_35mm"`
	ExposureBias    string  `json:"exposure_bias"`
	FlashFired      bool    `json:"flash_fired"`

	// location
	Lat      *float64 `json:"lat"`
	Lng      *float64 `json:"lng"`
	Altitude *float32 `json:"altitude"`

	// everything imagemeta decoded: IFD0, ExifIFD, GPS, IFD1/2, SubIFDs, MakerNote
	Exif datatypes.JSON `gorm:"type:jsonb" json:"exif"`

	CreatedAt time.Time `json:"created_at"`
}

func ExtractMetadata(data []byte) (Image, error) {
	m := Image{Size: len(data)}

	ex, exErr := imagemeta.Decode(bytes.NewReader(data))
	if exErr == nil {
		fillFromExif(&m, ex)
	}

	if cfg, format, err := image.DecodeConfig(bytes.NewReader(data)); err == nil {
		m.Format, m.Width, m.Height = format, cfg.Width, cfg.Height
	} else if exErr != nil {
		return Image{}, err
	}

	return m, nil
}

func fillFromExif(m *Image, ex exif.Exif) {
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
