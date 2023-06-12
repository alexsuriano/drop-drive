package folders

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/alexsuriano/drop-drive/internal/files"
)

func (h *handler) List(rw http.ResponseWriter, r *http.Request) {

	content, err := SelectRootFolderContent(h.db)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	folderContent := FolderContent{
		Folder: Folder{
			Name: "root",
		},
		Content: content,
	}

	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(folderContent)
}

func SelectRootFolderContent(db *sql.DB) ([]FolderResource, error) {
	subFolders, err := selectRootSubFolders(db)
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

	folderFiles, err := files.ListRoot(db)
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

func selectRootSubFolders(db *sql.DB) ([]Folder, error) {
	statement := `SELECT * FROM "folders" WHERE "parent_id"=$1 IS NULL AND "deleted"=false`

	rows, err := db.Query(statement)
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
