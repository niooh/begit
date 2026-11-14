package commands

import "begit/internal/cli"

// RegisterAll installs every command into cli's table.
func RegisterAll() {
	cli.Register(cli.Command{
		Name:  "init",
		Usage: "init",
		Desc:  "Initialize a begit repository in the current directory",
		Run:   cmdInit,
	})
	cli.Register(cli.Command{
		Name:  "send",
		Usage: `send <message>`,
		Desc:  "Create a full snapshot of the current directory",
		Run:   cmdSend,
		VerboseHelp: map[string]string{
			`<message>`: "Commit message, which will be stored in info.json",
		},
	})
	cli.Register(cli.Command{
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
	cli.Register(cli.Command{
		Name:  "ls",
		Usage: "ls",
		Desc:  "List snapshots",
		Run:   cmdLS,
	})
	cli.Register(cli.Command{
		Name:  "multi-select",
		Usage: "multi-select",
		Desc:  "Pick a new-old pair, then diff the pair or remove every snapshot between them",
		Run:   cmdMultiSelect,
		VerboseHelp: map[string]string{
			"new-old": "input, e.g. 1-3 (newer-older)",
			"d":       "diff new and old",
			"x":       "permanently delete every snapshot between old and new, inclusive",
			"b":       "back to snapshot list",
		},
	})
	cli.Register(cli.Command{
		Name:  "view",
		Usage: "view <relative-path>",
		Desc:  "View the change history of a path across snapshots",
		Run:   cmdView,
		VerboseHelp: map[string]string{
			"<relative-path>": "Exact path relative to the repo root",
		},
	})
	cli.Register(cli.Command{
		Name:  "remove-all",
		Usage: "remove-all <relative-path>",
		Desc:  "Remove a relative path from every snapshot",
		Run:   cmdRemoveAll,
		VerboseHelp: map[string]string{
			"<relative-path>": "Exact path relative to the repo root, e.g. rubbish.txt or rubbish/rubbish.txt",
		},
	})
	cli.Register(cli.Command{
		Name:  "md",
		Usage: "md",
		Desc:  "Export the current directory as a Markdown report",
		Run:   cmdExportMD,
	})
	cli.Register(cli.Command{
		Name:  "typ",
		Usage: "typ",
		Desc:  "Export the current directory as a Typst report",
		Run:   cmdExportTyp,
	})
}
