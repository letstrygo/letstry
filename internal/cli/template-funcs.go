package cli

import (
	"sort"
	"strings"
	"text/template"

	"github.com/fatih/color"
)

func defaultTemplateFuncs(app CliApp) template.FuncMap {
	return template.FuncMap{
		"getCallerName": func() string {
			return app.Name
		},
		"wrap": wrap,
		"getMaxWidth": func() int {
			return app.Config.DescriptionMaxWidth
		},
		"longestStringLength": longestStringLength,
		"padEnd":              padEnd,
		"white":               white,
		"whiteCommand":        whiteCommand,
		"commands": func() []Command {
			return commands(app)
		},
	}
}

func commands(app CliApp) []Command {
	commands := make([]Command, 0, len(app.commands))
	for _, cmd := range app.commands {
		commands = append(commands, cmd)
	}

	if app.Config.HelpCommandSorter != nil {
		sort.Slice(commands, func(i int, j int) bool {
			return app.Config.HelpCommandSorter(commands[i], commands[j])
		})
	} else {
		sort.Slice(commands, func(i int, j int) bool {
			return CommandSorterAlphabeticalFollowedByHelp(commands[i], commands[j])
		})
	}

	return commands
}

func wrap(text string, lineLength int, lineStartPadding int) string {
	if lineLength < 1 {
		return text
	}

	padding := strings.Repeat(" ", lineStartPadding)
	var result strings.Builder
	words := strings.Fields(text)

	currentLineLength := lineStartPadding

	for i, word := range words {
		if currentLineLength+len(word)+1 > lineLength {
			result.WriteString("\n" + padding)
			currentLineLength = lineStartPadding
		} else if i > 0 {
			result.WriteString(" ")
			currentLineLength++
		}
		result.WriteString(word)
		currentLineLength += len(word)
	}

	return result.String()
}

func longestStringLength(cmds []Command) int {
	var max int
	for _, s := range cmds {
		if len(s.Name) > max {
			max = len(s.Name)
		}
	}
	return max
}

func padEnd(val string, length int) string {
	if length <= len(val) {
		return val
	}
	return val + strings.Repeat(" ", length-len(val))
}

func white(s string) string {
	return color.HiWhiteString(s)
}

func whiteCommand(cmd string) string {
	return color.HiWhiteString(cmd)
}
