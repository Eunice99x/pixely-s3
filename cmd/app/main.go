package main

import (
	"context"
	"log"
	"os"

	"github.com/eunice99x/pixely/kyc/internal/database"
	"github.com/eunice99x/pixely/kyc/internal/handlers"
	"github.com/eunice99x/pixely/kyc/internal/storage"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	db, err := database.New(getenv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5433/ecomm?sslmode=disable"))
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}

	store, err := storage.NewS3(context.Background())
	if err != nil {
		log.Fatalf("error creating s3: %v", err)
	}

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	h := handlers.New(db, store)
	e.POST("/upload", h.Upload)
	e.GET("/images/:id", h.GetMetadata)
	e.GET("/health", h.Health)

	if err := e.Start(":" + getenv("PORT", "3000")); err != nil {
		log.Fatal(err)
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
