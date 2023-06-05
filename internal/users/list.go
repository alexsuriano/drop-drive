package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *handler) list(rw http.ResponseWriter, r *http.Request) {
	userList, err := SelectAll(h.db)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(userList)
}

func SelectAll(db *sql.DB) ([]User, error) {
	statement := `SELECT * FROM "user" WHERE "deleted"=false`
	rows, err := db.Query(statement)

	userList := make([]User, 0)
	for rows.Next() {
		u := User{}

		err = rows.Scan(&u.ID, &u.Name, &u.Password, &u.CreatedAt,
			&u.ModifiedAt, &u.Deleted, &u.LastLogin)
		if err != nil {
			continue
			// return nil, err
		}

		userList = append(userList, u)
	}

	return userList, nil
}
