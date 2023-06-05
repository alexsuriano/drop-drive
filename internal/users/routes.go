package users

import (
	"database/sql"

	"github.com/go-chi/chi"
)

type handler struct {
	db *sql.DB
}

func SetRoutes(r chi.Router, db *sql.DB) {
	handler := handler{
		db: db,
	}

	r.Post("/", handler.Create)
	r.Put("/{id}", handler.Modify)
	r.Delete("/{id}", handler.Delete)
	r.Get("/{id}", handler.GetByID)
	r.Get("/", handler.list)
}
