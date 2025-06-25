package manager

import (
	"context"

	"github.com/go-git/go-git/v5"
)

type UpdateTemplateArguments struct {
	TemplateName string
}

func (m *manager) UpdateTemplate(ctx context.Context, arg UpdateTemplateArguments) error {
	t, err := m.GetTemplate(ctx, arg.TemplateName)
	if err != nil {
		return err
	}

	absPath := t.AbsolutePath(ctx)

	repo, err := git.PlainOpen(absPath)
	if err != nil {
		return err
	}

	err = repo.Fetch(&git.FetchOptions{
		Force: true,
	})
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		Force: true,
	})
	if err != nil {
		return err
	}

	return nil
}
