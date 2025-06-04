package manager

import (
	"context"
	"fmt"
	"os"

	"github.com/letstrygo/letstry/internal/util/identifier"
)

var (
	ErrCurrentDirectoryIsNotASession = fmt.Errorf("current directory is not a session")
	ErrSessionNotFound               = fmt.Errorf("session not found")
)

// GetSession returns the session with the given ID
func (s *manager) GetSession(ctx context.Context, id identifier.ID) (Session, error) {
	sessions, err := s.ListSessions(ctx)
	if err != nil {
		return Session{}, err
	}

	for _, session := range sessions {
		if session.ID == id {
			return session, nil
		}
	}

	return Session{}, fmt.Errorf("session with ID %s not found", id)
}

// GetCurrentSession returns the session for the current working directory
func (s *manager) GetCurrentSession(ctx context.Context) (Session, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Session{}, fmt.Errorf("failed to get current working directory: %v", err)
	}

	sess, err := s.GetSessionForPath(ctx, cwd)
	if err != nil && err == ErrSessionNotFound {
		return sess, ErrCurrentDirectoryIsNotASession
	}

	return sess, err
}

// GetSessionForPath returns the session for the given path
func (s *manager) GetSessionForPath(ctx context.Context, path string) (Session, error) {
	return s.GetSessionForPredicate(ctx, func(sess Session) bool {
		return sess.Location == path
	})
}

// GetSessionForPredicate returns the session that matches the given predicate
func (s *manager) GetSessionForPredicate(ctx context.Context, predicate func(Session) bool) (Session, error) {
	// get the list of sessions
	sessions, err := s.ListSessions(ctx)
	if err != nil {
		return Session{}, err
	}

	// find the session with the same location
	for _, session := range sessions {
		if predicate(session) {
			return session, nil
		}
	}

	return Session{}, ErrSessionNotFound
}
