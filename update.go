package main

var cmdUpdate = &Command{
	Run:      runUpdate,
	Usage:    "update",
	Category: "mc",
	Long: `
Update downloads and installs the next version of mc.

This command is unlisted as users never have to run it directly.
`,
}

func runUpdate(cmd *Command, args []string) {
	if updater == nil {
		printFatal("Dev builds don't support auto-updates")
	}
	// do update
}

type Updater struct{}

func (u *Updater) backgroundRun() {}
