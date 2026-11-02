package model

import (
	"time"

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
