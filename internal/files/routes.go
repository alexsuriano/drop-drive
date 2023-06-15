package files

import (
	"database/sql"

	"github.com/alexsuriano/drop-drive/internal/bucket"
	"github.com/alexsuriano/drop-drive/internal/queue"
	"github.com/go-chi/chi"
)

type handler struct {
	db     *sql.DB
	bucket *bucket.Bucket
	queue  *queue.Queue
}

func SetRoutes(r chi.Router, db *sql.DB, b *bucket.Bucket, q *queue.Queue) {
	handler := handler{
		db:     db,
		bucket: b,
		queue:  q,
	}

	r.Post("/", handler.Create)
	r.Put("/{id}", handler.Modify)
	r.Get("/{id}", handler.Get)
	r.Delete("/{id}", handler.Delete)

}
