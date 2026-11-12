package cli

type Command struct {
	Name        string
	Usage       string
	Desc        string
	Run         func(args []string) int
	VerboseHelp map[string]string
}

var commands = map[string]Command{}

func register(c Command) {
	commands[c.Name] = c
}

func init() {
	register(Command{
		Name:  "init",
		Usage: "init",
		Desc:  "Initialize a begit repository in the current directory",
		Run:   cmdInit,
	})
	register(Command{
		Name:  "send",
		Usage: `send "message"`,
		Desc:  "Create a full snapshot of the current directory",
		Run:   cmdSend,
		VerboseHelp: map[string]string{
			`"message"`: "Commit message stored in info.json",
		},
	})
	register(Command{
		Name:  "select",
		Usage: "select",
		Desc:  "Interactively select and operate on snapshots",
		Run:   cmdSelect,
		VerboseHelp: map[string]string{
			"d": "diff selected snapshot against the previous one",
			"r": "reset working directory to selected snapshot",
			"x": "permanently delete selected snapshot",
			"b": "back to snapshot list",
		},
	})
	register(Command{
		Name:  "remove-all",
		Usage: "remove-all <relative-path>",
		Desc:  "Remove a relative path from every snapshot",
		Run:   cmdRemoveAll,
		VerboseHelp: map[string]string{
			"<relative-path>": "Exact path relative to the repo root, e.g. rubbish.txt or rubbish/rubbish.txt",
		},
	})
	register(Command{
		Name:  "md",
		Usage: "md",
		Desc:  "Export the current directory as a Markdown report",
		Run:   cmdExportMD,
	})
	register(Command{
		Name:  "typ",
		Usage: "typ",
		Desc:  "Export the current directory as a Typst report",
		Run:   cmdExportTyp,
	})
}
