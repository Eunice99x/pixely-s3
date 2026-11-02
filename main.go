package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/eunice99x/pixely/kyc/db"
	"github.com/eunice99x/pixely/kyc/handler"
	"github.com/eunice99x/pixely/kyc/storage"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	godotenv.Load()

	gdb, err := db.NewDatabase()
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

	h := handler.New(gdb, store)
	e.POST("/upload", h.Upload)
	e.GET("/images/:id", h.GetMetadata)

	e.GET("/health", func(c *echo.Context) error {
		return c.String(http.StatusOK, "suceess!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}
}
