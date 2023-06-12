package folders

import (
	"errors"
	"time"
)

var (
	ErrNameRequired = errors.New("name is required")
)

type Folder struct {
	ID         int64     `json:"id"`
	ParentID   int64     `json:"parent_id"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
	Deleted    bool      `json:"-"`
}

type FolderContent struct {
	Folder  Folder           `json:"folder"`
	Content []FolderResource `json:"content"`
}

type FolderResource struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}

func New(name string, parentID int64) (*Folder, error) {
	folder := Folder{
		ParentID: parentID,
		Name:     name,
	}

	err := folder.Validate()
	if err != nil {
		return nil, err
	}

	return &folder, nil
}

func (f *Folder) Validate() error {
	if f.Name == "" {
		return ErrNameRequired
	}

	return nil
}
