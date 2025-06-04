package manager

import (
	"context"
	"fmt"

	"github.com/letstrygo/letstry/internal/storage"
)

var (
	ErrManagerNotFound = fmt.Errorf("manager not found in context")
)

type mgrCtxKey struct {
	value string
}

var (
	mgrKey = mgrCtxKey{"manager"}
)

type manager struct {
	storage *storage.Storage
}

// ContextWithManager returns a new context with the session manager
func ContextWithManager(ctx context.Context) context.Context {
	return context.WithValue(ctx, mgrKey, manager{
		storage: storage.GetStorage(),
	})
}

// GetManager returns the session manager
func GetManager(ctx context.Context) (*manager, error) {
	mgr := ctx.Value(mgrKey)
	if mgr != nil {
		if mgr, ok := mgr.(manager); ok {
			return &mgr, nil
		}
	}

	return nil, ErrManagerNotFound
}
