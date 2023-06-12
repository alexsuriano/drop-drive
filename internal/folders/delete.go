package folders

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/alexsuriano/drop-drive/internal/files"
	"github.com/go-chi/chi"
)

func (h *handler) Delete(rw http.ResponseWriter, r *http.Request) {
	folderID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = deleteFolderContent(h.db, int64(folderID))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	err = Delete(h.db, int64(folderID))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
}

func Delete(db *sql.DB, idFolder int64) error {

	statement := `UPDATE "folders" SET "modified_at"=$1 "deleted"=true WHERE "id"=$2`

	_, err := db.Exec(statement, time.Now(), idFolder)

	return err
}

func deleteFolderContent(db *sql.DB, folderID int64) error {
	err := deleteFiles(db, folderID)
	if err != nil {
		return err
	}

	return deleteSubFolders(db, folderID)
}

func deleteFiles(db *sql.DB, folderID int64) error {
	fileList, err := files.List(db, folderID)
	if err != nil {
		return err
	}

	removedFiles := make([]files.File, 0, len(fileList))

	for _, file := range fileList {
		file.Deleted = true
		err := files.Update(db, &file, file.ID)
		if err != nil {
			break
		}

		removedFiles = append(removedFiles, file)
	}

	if len(fileList) != len(removedFiles) {
		for _, file := range removedFiles {
			file.Deleted = false
			files.Update(db, &file, file.ID)
		}
	}

	return nil
}

func deleteSubFolders(db *sql.DB, folderID int64) error {
	subFoldersList, err := selectSubFolders(db, folderID)
	if err != nil {
		return err
	}

	removedFolders := make([]Folder, 0, len(subFoldersList))

	for _, subFolder := range subFoldersList {
		err := Delete(db, subFolder.ID)
		if err != nil {
			break
		}

		err = deleteFolderContent(db, subFolder.ID)
		if err != nil {
			Update(db, &subFolder, subFolder.ID)
			break
		}

		removedFolders = append(removedFolders, subFolder)
	}

	if len(subFoldersList) != len(removedFolders) {
		for _, subFolder := range removedFolders {
			subFolder.Deleted = false
			Update(db, &subFolder, subFolder.ID)
			// err := Update(db, &subFolder, subFolder.ID)
			// if err != nil {
			// 	return err
			// }

		}
	}

	return nil
}
