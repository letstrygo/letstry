package manager

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/letstrygo/letstry/internal/config/editors"
	"github.com/letstrygo/letstry/internal/util/access"
	"github.com/letstrygo/letstry/internal/util/identifier"
)

type Source struct {
	SourceType SessionSourceType `json:"sourceType"`
	Value      string            `json:"value"`
}

func (s Source) String() string {
	switch s.SourceType {
	case SessionSourceTypeBlank:
		return fmt.Sprintf("[%s]", s.SourceType)
	default:
		return fmt.Sprintf("[%s, %s]", s.SourceType, s.Value)
	}
}

func (s Source) FormattedValue() string {
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
	Source   Source         `json:"source"`
	Editor   editors.Editor `json:"editor"`
}

func (s *Session) IsActive() bool {
	return access.IsPathUse(s.Location)
}

func (s *Session) String() string {
	src := s.Source.FormattedValue()
	id := s.ID.FormattedString()
	editor := color.BlueString("(%s)", s.Editor.Name)

	return fmt.Sprintf("id=%s, editor=%s, src=%s", id, editor, src)
}
