package manager

import (
	"context"
	"fmt"

	"github.com/letstrygo/letstry/internal/logging"
	"github.com/letstrygo/letstry/internal/util/identifier"
)

type PruneSessionArguments struct {
	SessionID identifier.ID
}

func (s *manager) PruneSession(ctx context.Context, args PruneSessionArguments) error {
	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return err
	}

	session, err := s.GetSession(ctx, args.SessionID)
	if err != nil {
		return err
	}

	if session.IsActive() {
		return fmt.Errorf("cannot prune session: %s (directory still being accessed)", session.ID.FormattedString())
	}

	logger.Printf("pruning inactive session: %s\n", session.ID.FormattedString())
	return s.removeSession(ctx, session.ID)
}
