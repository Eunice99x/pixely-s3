package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/eunice99x/pixely/kyc/db"
	"github.com/eunice99x/pixely/kyc/handler"
	"github.com/eunice99x/pixely/kyc/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}
	defer db.Close()

	store, err := storage.NewS3(context.Background())
	if err != nil {
		log.Fatalf("error creating s3: %v", err)
	}

	h := handler.New(db.GetDB(), store)
	r.Post("/upload", h.Upload)
	r.Get("/images/{id}", h.GetMetadata)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("suceess!"))
	})

	fmt.Println("Server starting on :3000...")
	if err := http.ListenAndServe(":3000", r); err != nil {
		panic(err)
	}
}
