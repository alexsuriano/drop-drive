package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

func (h *handler) Modify(rw http.ResponseWriter, r *http.Request) {
	user := &User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.Name == "" {
		http.Error(rw, ErrNameRequired.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	err = Update(h.db, user, int64(id))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	//TODO Get ID

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(user)
}

func Update(db *sql.DB, user *User, id int64) error {
	user.ModifiedAt = time.Now()

	statement := `UPDATE "users" SET "name"=$1, "modified_at"=$2 where "id"=$3`

	_, err := db.Exec(statement, user.Name, user.ModifiedAt, id)

	return err
}
