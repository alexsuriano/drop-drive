package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) Create(rw http.ResponseWriter, r *http.Request) {
	user := &User{}

	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = user.SetPassword(user.Password)
	if err != nil {

	}

	err = user.Validate()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := Insert(h.db, user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	user.ID = id

	rw.Header().Add("Content-Type", "applicantion/json")
	json.NewEncoder(rw).Encode(user)

}

func Insert(db *sql.DB, u *User) (int64, error) {
	statement := `INSERT INTO "users" ("name", "login", "password", "modified_at") VALUES ($1, $2, $3, $4)`

	result, err := db.Exec(statement, u.Name, u.Login, u.Password, u.ModifiedAt, u.Deleted)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
