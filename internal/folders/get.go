package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/alexsuriano/drop-drive/internal/files"
	"github.com/go-chi/chi"
)

func (h *handler) Get(rw http.ResponseWriter, r *http.Request) {

	folderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	folder, err := SelectFolder(h.db, int64(folderID))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	content, err := SelectFolderContent(h.db, int64(folderID))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	folderContent := FolderContent{
		Folder:  *folder,
		Content: content,
	}

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(folderContent)

}

func SelectFolder(db *sql.DB, folderID int64) (*Folder, error) {
	statement := `SELECT * FROM "folders" WHERE "id"=$1`
	row := db.QueryRow(statement, folderID)

	var f Folder
	err := row.Scan(&f.ID, &f.ParentID, &f.Name,
		&f.CreatedAt, &f.ModifiedAt, &f.Deleted)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func selectSubFolders(db *sql.DB, folderID int64) ([]Folder, error) {
	statement := `SELECT * FROM "folders" WHERE "parent_id"=$1 and "deleted"=false`

	rows, err := db.Query(statement, folderID)
	if err != nil {
		return nil, err
	}

	subFolders := make([]Folder, 0)
	for rows.Next() {
		var folder Folder
		err := rows.Scan(&folder.ID, &folder.ParentID, &folder.Name,
			&folder.CreatedAt, &folder.ModifiedAt, &folder.Deleted)
		if err != nil {
			continue
		}

		subFolders = append(subFolders, folder)
	}

	return subFolders, nil
}

func SelectFolderContent(db *sql.DB, folderID int64) ([]FolderResource, error) {
	subFolders, err := selectSubFolders(db, folderID)
	if err != nil {
		return nil, err
	}

	folderResource := make([]FolderResource, 0, len(subFolders))
	for _, subFolder := range subFolders {
		resource := FolderResource{
			ID:         subFolder.ID,
			Name:       subFolder.Name,
			Type:       "directory",
			CreatedAt:  subFolder.CreatedAt,
			ModifiedAt: subFolder.ModifiedAt,
		}

		folderResource = append(folderResource, resource)
	}

	folderFiles, err := files.List(db, folderID)
	if err != nil {
		return nil, err
	}

	for _, file := range folderFiles {
		resource := FolderResource{
			ID:         file.ID,
			Name:       file.Name,
			Type:       file.Type,
			CreatedAt:  file.CreatedAt,
			ModifiedAt: file.ModifiedAt,
		}

		folderResource = append(folderResource, resource)
	}

	return folderResource, nil
}
