package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

func (h *handler) Modify(rw http.ResponseWriter, r *http.Request) {
	folder := new(Folder)

	err := json.NewDecoder(r.Body).Decode(folder)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = folder.Validate()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Update(h.db, folder, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	//Todo GET ID

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(folder)
}

func Update(db *sql.DB, f *Folder, id int64) error {
	f.ModifiedAt = time.Now()

	statement := `UPDATE "folder" SET "name"=$1, "modified_at"=$2 WHERE "id"=$3`

	_, err := db.Exec(statement, f.Name, f.ModifiedAt, id)

	return err
}
