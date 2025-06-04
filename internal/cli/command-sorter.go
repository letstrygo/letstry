package cli

type CommandSorter func(a Command, b Command) bool

func CommandSorterAlphabetical(a Command, b Command) bool {
	return a.Name < b.Name
}

func CommandSorterAlphabeticalFollowedByHelp(a Command, b Command) bool {
	if a.Name == "help" {
		return false
	}

	if b.Name == "help" {
		return true
	}

	return CommandSorterAlphabetical(a, b)
}

func CommandSorterOrderedAs(commands []Command) CommandSorter {
	return func(a, b Command) bool {
		// order them in the same order as the commands slice
		for _, command := range commands {
			if command.Name == a.Name {
				return true
			}
			if command.Name == b.Name {
				return false
			}
		}

		return false
	}
}
