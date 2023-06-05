package folders

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
}
