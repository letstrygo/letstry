package repository

import (
	"github.com/go-git/go-git/v6"
)

const (
	RepositoryURL string = "https://github.com/letstrygo/templates.git"
)

type Repository struct {
	location      string
	gitRepository *git.Repository
}

func NewRepository(location string) *Repository {
	return &Repository{location, nil}
}

func (r *Repository) Init() error {
	var err error

	r.gitRepository, err = git.PlainOpen(r.location)
	if err != nil {
		if err != git.ErrRepositoryNotExists {
			return err
		}

		r.gitRepository, err = git.PlainClone(r.location, &git.CloneOptions{
			URL: RepositoryURL,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
