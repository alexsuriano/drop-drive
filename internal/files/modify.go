package files

import (
	"database/sql"
	"net/http"
	"time"
)

func Modify(rw http.ResponseWriter, r http.Request) {

}

func Update(db *sql.DB, f *File, fileID int64) error {
	f.ModifiedAt = time.Now()

	statement := `UPDATE "files" SET "name"=$1, "modified_at"=$2, "deleted"=$3 WHERE "id"=$4`
	_, err := db.Exec(statement, f.Name, f.ModifiedAt, f.Deleted, fileID)

	return err
}
