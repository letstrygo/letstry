package manager

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/letstrygo/letstry/internal/config/editors"
	"github.com/letstrygo/letstry/internal/util/access"
	"github.com/letstrygo/letstry/internal/util/identifier"
)

type sessionSource struct {
	SourceType SessionSourceType `json:"sourceType"`
	Value      string            `json:"value"`
}

func (s sessionSource) String() string {
	switch s.SourceType {
	case SessionSourceTypeBlank:
		return fmt.Sprintf("[%s]", s.SourceType)
	default:
		return fmt.Sprintf("[%s, %s]", s.SourceType, s.Value)
	}
}

func (s sessionSource) FormattedString() string {
	var colorWrapper func(format string, a ...interface{}) string = color.WhiteString

	switch s.SourceType {
	case SessionSourceTypeDirectory:
		fallthrough
	case SessionSourceTypeRepository:
		colorWrapper = color.HiBlueString
	case SessionSourceTypeTemplate:
		colorWrapper = color.HiMagentaString
	case SessionSourceTypeBlank:
		colorWrapper = color.HiWhiteString
	}

	return colorWrapper("%s", s.String())
}

type Session struct {
	ID       identifier.ID  `json:"id"`
	Location string         `json:"location"`
	PID      int            `json:"pid"`
	Source   sessionSource  `json:"source"`
	Editor   editors.Editor `json:"editor"`
}

func (s *Session) IsActive() bool {
	return access.IsPathUse(s.Location)
}

func (s *Session) String() string {
	src := s.Source.FormattedString()
	id := s.FormattedID()
	editor := color.BlueString("(%s)", s.Editor.Name)

	return fmt.Sprintf("id=%s, editor=%s, src=%s", id, editor, src)
}

func (s *Session) FormattedID() string {
	return color.HiGreenString(s.ID.String())
}
