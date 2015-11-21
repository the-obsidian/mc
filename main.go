package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mgutz/ansi"
	"github.com/the-obsidian/mc/rollbar"
	"github.com/the-obsidian/mc/term"
)

var commands = []*Command{
	cmdInstall,
	//cmdVersion,
	cmdHelp,

	//helpCommands,
	//helpAbout,

	// unlisted
	//cmdUpdate,
}

func main() {
	log.SetFlags(0)

	args := os.Args[1:]
	if len(args) < 1 || strings.IndexRune(args[0], '-') == 0 {
		printUsageTo(os.Stderr)
		os.Exit(2)
	}

	// Perform updates as early as possible
	if args[0] == cmdUpdate.Name() {
		cmdUpdate.Run(cmdUpdate, args)
		return
	} else if updater != nil {
		defer updater.backgroundRun()
	}

	if !term.IsANSI(os.Stdout) {
		ansi.DisableColors(true)
	}

	for _, cmd := range commands {
		if cmd.Name() != args[0] || cmd.Run == nil {
			continue
		}

		defer recoverPanic()

		cmd.Flag.SetDisableDuplicates(true) // allow duplicate flag options
		cmd.Flag.Usage = func() {
			cmd.PrintUsage()
		}

		if cmd.NeedsServer {
			if exists, err := fileExists("server.yml"); err != nil || !exists {
				printFatal("server.yml not found - is this a server directory?")
			}
		}

		if err := cmd.Flag.Parse(args[1:]); err == flag.ErrHelp {
			cmdHelp.Run(cmdHelp, args[:1])
		} else if err != nil {
			printError(err.Error())
			os.Exit(2)
		}
		cmd.Run(cmd, cmd.Flag.Args())
		return
	}

	fmt.Fprintf(os.Stderr, "Unknown command: %s\n", args[0])
	if g := suggest(args[0]); len(g) > 0 {
		fmt.Fprintf(os.Stderr, "Possible alternatives: %v\n", strings.Join(g, " "))
	}
	fmt.Fprintf(os.Stderr, "Run 'mc help' for usage.\n")
	os.Exit(2)
}

var rollbarClient = &rollbar.Client{
	AppName:    "mc",
	AppVersion: Version,
	Endpoint:   "https://api.rollbar.com/api/1/item/",
	Token:      "c0fd71dc4a724934b0826c73ef3cd269",
}

func recoverPanic() {
	if Version == "dev" {
		return
	}

	if rec := recover(); rec != nil {
		message := ""
		switch rec := rec.(type) {
		case error:
			message = rec.Error()
		default:
			message = fmt.Sprintf("%v", rec)
		}
		if err := rollbarClient.Report(message); err != nil {
			printError("reporting crash failed: %s", err.Error())
			panic(err)
		}
		printFatal("mc encountered and reported an internal client issue")
	}
}
