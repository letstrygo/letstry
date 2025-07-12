package repository

import (
	"io"
	"net/http"
	"os"

	"github.com/letstrygo/letstry/internal/storage"
	"github.com/letstrygo/templates"
)

const (
	RepositoryURL string = "https://raw.githubusercontent.com/letstrygo/templates/refs/heads/main/dist/database.sqlite"
)

type Repository struct {
	*templates.Connection
	location string
}

func NewRepository() (*Repository, error) {
	store := storage.GetStorage()
	location := store.GetAbsolutePath("repository.db")

	repository := &Repository{nil, location}
	if err := repository.init(); err != nil {
		return nil, err
	}

	return repository, nil
}

func (r *Repository) init() error {
	_, err := os.Stat(r.location)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		// Download the remote repository.

		resp, err := http.Get(RepositoryURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		out, err := os.Create(r.location)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return err
		}
	}

	r.Connection, err = templates.NewConnectionWithPath(r.location)
	if err != nil {
		return err
	}

	return nil
}
