package files

import "database/sql"

func List(db *sql.DB, folderID int64) ([]File, error) {
	statement := `SELECT * FROM "files" WHERE "folder_id"=$1`

	rows, err := db.Query(statement, folderID)
	if err != nil {
		return nil, err
	}

	fileList := make([]File, 0)
	for rows.Next() {
		f := File{}

		err := rows.Scan(&f.ID, &f.FolderID, &f.OwnerID, &f.Name,
			&f.Type, &f.Path, &f.CreatedAt, &f.ModifiedAt, &f.Deleted)
		if err != nil {
			continue
		}

		fileList = append(fileList, f)

	}

	return fileList, nil
}
