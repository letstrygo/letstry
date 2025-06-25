package manager

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/letstrygo/letstry/internal/logging"
	"github.com/letstrygo/letstry/internal/util/identifier"
)

var (
	ErrInvalidSessionDisplayType = errors.New("invalid session display type")
)

type SessionDisplayType string

const (
	SessionDisplayTypeFull     SessionDisplayType = "full"
	SessionDisplayTypeLocation SessionDisplayType = "location"
	SessionDisplayTypePID      SessionDisplayType = "pid"
	SessionDisplayTypeEditor   SessionDisplayType = "editor"
	SessionDisplayTypeJSON     SessionDisplayType = "json"
)

var (
	SessionDisplayTypes []SessionDisplayType = []SessionDisplayType{
		SessionDisplayTypeFull,
		SessionDisplayTypeLocation,
		SessionDisplayTypePID,
		SessionDisplayTypeEditor,
		SessionDisplayTypeJSON,
	}
)

func ParseSessionDisplayType(v string) (SessionDisplayType, error) {
	for _, dt := range SessionDisplayTypes {
		if string(dt) == v {
			return dt, nil
		}
	}

	return SessionDisplayTypeFull, ErrInvalidSessionDisplayType
}

type ShowSessionArguments struct {
	SessionID  *identifier.ID
	OutputType SessionDisplayType
}

func (m *manager) ShowSession(ctx context.Context, args ShowSessionArguments) error {
	var (
		session Session
		err     error
	)

	switch {
	case args.SessionID != nil:
		// Try to load session from the ID.
		session, err = m.GetSession(ctx, *args.SessionID)
	default:
		session, err = m.GetCurrentSession(ctx)
	}
	if err != nil {
		return err
	}

	log, err := logging.LoggerFromContext(ctx)
	if err != nil {
		return err
	}

	switch args.OutputType {
	case SessionDisplayTypeFull:
		log.Printf("%v", session.String())
	case SessionDisplayTypeLocation:
		println(session.Location)
	case SessionDisplayTypePID:
		println(session.PID)
	case SessionDisplayTypeEditor:
		log.Printf("%v", session.Editor.FullString())
	case SessionDisplayTypeJSON:
		sessionData, err := json.MarshalIndent(session, "", "    ")
		if err != nil {
			return err
		}

		println(string(sessionData))
	default:
		return ErrInvalidSessionDisplayType
	}

	return nil
}
