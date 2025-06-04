package environment

import (
	"context"
	"errors"
	"os"
)

type environmentCtxKey struct {
	Name string
}

var (
	environmentKey environmentCtxKey = environmentCtxKey{"environment"}
)

var (
	ErrEnvironmentNotFound = errors.New("environment not found")
)

type Environment struct {
	DebuggerAttached bool
}

func EnvironmentFromContext(ctx context.Context) (Environment, error) {
	if env, ok := ctx.Value(environmentKey).(Environment); ok {
		return env, nil
	}

	return Environment{}, ErrEnvironmentNotFound
}

func ContextWithEnvironment(ctx context.Context) context.Context {
	return context.WithValue(ctx, environmentKey, Environment{
		DebuggerAttached: os.Getenv("DEBUGGER_ATTACHED") == "true",
	})
}
