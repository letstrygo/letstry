package manager

import (
	"context"
	"fmt"

	"github.com/letstrygo/letstry/internal/logging"
	"github.com/letstrygo/letstry/internal/util/identifier"
)

type CleanSessionArguments struct {
	SessionID identifier.ID
}

func (s *manager) CleanSession(ctx context.Context, args CleanSessionArguments) error {
	logger, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return err
	}

	session, err := s.GetSession(ctx, args.SessionID)
	if err != nil {
		return err
	}

	if session.IsActive() {
		return fmt.Errorf("cannot clean session: %s (directory still being accessed)", session.FormattedID())
	}

	logger.Printf("cleaning inactive session: %s\n", session.FormattedID())
	return s.removeSession(ctx, session.ID)
}
