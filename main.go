package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eunice99x/pixely/kyc/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}
	defer db.Close()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("suceess!"))
	})

	fmt.Println("Server starting on :3000...")
	if err := http.ListenAndServe(":3000", r); err != nil {
		panic(err)
	}
}
