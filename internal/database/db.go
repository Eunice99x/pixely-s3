package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/eunice99x/pixely/kyc/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ErrNotFound = errors.New("not found")

type DB struct {
	db *gorm.DB
}

func New(dsn string) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	return &DB{db: db}, nil
}

func (d *DB) CreateImage(ctx context.Context, img *model.Image) error {
	return d.db.WithContext(ctx).Create(img).Error
}

func (d *DB) GetImage(ctx context.Context, id uint) (model.Image, error) {
	var img model.Image
	err := d.db.WithContext(ctx).First(&img, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return img, ErrNotFound
	}
	return img, err
}

func (d *DB) Ping(ctx context.Context) error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
