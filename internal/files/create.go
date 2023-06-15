package files

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alexsuriano/drop-drive/internal/queue"
)

func (h *handler) Create(rw http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	path := fmt.Sprintf("/%s", fileHeader.Filename)

	err = h.bucket.Upload(file, path)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	entity, err := New(1, fileHeader.Filename, fileHeader.Header.Get("Content-Type"), path)
	if err != nil {
		h.bucket.Delete(path)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	folderID := r.Form.Get("folder_id")
	if folderID != "" {
		fID, err := strconv.Atoi(folderID)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		entity.FolderID = int64(fID)
	}

	id, err := Insert(h.db, entity)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	entity.ID = id
	dto := queue.QueueDTO{
		Filename: fileHeader.Filename,
		Path:     path,
		ID:       int(id),
	}

	msg, err := dto.Marshal()
	if err != nil {
		//TODO: rollback
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.queue.Publish(msg)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(entity)

}

func Insert(db *sql.DB, file *File) (int64, error) {
	statement := `INSERT INTO "Files" ("folder_id", "owner_id", "name", "type", 
	"path", "modified_at") VALUES ($1, $2, $3, $4, $5, $6)`

	result, err := db.Exec(statement, file.FolderID, file.OwnerID, file.Name,
		file.Type, file.Path, file.ModifiedAt)
	if err != nil {
		return -1, err
	}

	return result.LastInsertId()
}
