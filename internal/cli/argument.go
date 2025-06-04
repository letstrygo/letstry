package cli

import (
	"fmt"

	"github.com/fatih/color"
)

type Argument struct {
	Name        string
	Description string
	Required    bool
}

func (a Argument) Label() string {
	whiteName := color.HiWhiteString(a.Name)

	if !a.Required {
		return fmt.Sprintf("%s %s", whiteName, color.BlueString("(optional)"))
	}

	return fmt.Sprintf("%s %s", whiteName, color.RedString("(required)"))
}
