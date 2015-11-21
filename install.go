package main

import "os"

var cmdInstall = &Command{
	NeedsServer: true,
	Usage:       "install",
	Category:    "app",
	Short:       "install dependencies",
	Long: `
Install installs any required dependencies.

Examples:

	$ mc install
`,
}

func init() {
	cmdInstall.Run = runInstall
}

func runInstall(cmd *Command, args []string) {
	if len(args) != 0 {
		printUsageTo(os.Stderr)
		os.Exit(2)
	}

	config, err := NewConfigFromFile("server.yml")
	if err != nil {
		printFatal("Error reading config file: %v", err)
	}

	err = config.InstallPlugins()
	if err != nil {
		printFatal("failed to install plugins: %v", err)
	}
}
