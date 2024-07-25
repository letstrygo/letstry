package session_manager

import (
	"context"
	"encoding/json"
	"fmt"
)

func (s *sessionManager) ListSessions(ctx context.Context) ([]session, error) {
	var sessions []session = make([]session, 0)

	var defaultSessions []byte
	defaultSessions, err := json.MarshalIndent(sessions, "", "    ")
	if err != nil {
		return sessions, fmt.Errorf("failed to marshal default sessions: %v", err)
	}

	file, err := s.storage.OpenFileWithDefaultContent("sessions.json", defaultSessions)
	if err != nil {
		return sessions, fmt.Errorf("failed to open sessions file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&sessions)
	if err != nil {
		return sessions, fmt.Errorf("failed to decode sessions file: %v", err)
	}

	return sessions, nil
}
