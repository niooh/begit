package cli

type Command struct {
	Name        string
	Usage       string
	Desc        string
	Run         func(args []string) int
	VerboseHelp map[string]string
}

var commands = map[string]Command{}

func Register(c Command) {
	commands[c.Name] = c
}
