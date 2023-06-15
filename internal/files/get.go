package files

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (h *handler) Get(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var file *File
	file, err = Select(h.db, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(file)
}

func Select(db *sql.DB, id int64) (*File, error) {
	statement := `SELECT * FROM "files" WHERE "id"=$1`
	row := db.QueryRow(statement, id)

	var f File
	err := row.Scan(&f.ID, &f.FolderID, &f.OwnerID, &f.Name, &f.Type,
		&f.Path, &f.CreatedAt, &f.ModifiedAt, &f.Deleted)
	if err != nil {
		return nil, err
	}

	return &f, nil
}
