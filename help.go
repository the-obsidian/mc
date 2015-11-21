package main

import (
	"io"
	"log"
	"os"
	"text/template"
)

var cmdHelp = &Command{
	Usage:    "help [<topic>]",
	Category: "mc",
	Long:     `Help shows usage for a command or other topic.`,
}

func init() {
	cmdHelp.Run = runHelp // breaks init loop
}

func runHelp(cmd *Command, args []string) {
	if len(args) == 0 {
		printUsageTo(os.Stdout)
		return
	}

	if len(args) != 1 {
		printFatal("too many arguments")
	}

	switch args[0] {
	default:
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.PrintLongUsage()
			return
		}
	}

	log.Printf("Unknown help topic: %q. Run 'mc help'.\n", args[0])
	os.Exit(2)
}

func maxStrLen(strs []string) (strlen int) {
	for i := range strs {
		if len(strs[i]) > strlen {
			strlen = len(strs[i])
		}
	}
	return
}

var usageTemplate = template.Must(template.New("usage").Parse(`
Usage: mc <command> [options] [arguments]
Commands:
{{range .Commands}}{{if .Runnable}}{{if .List}}
    {{.Name | printf (print "%-" $.MaxRunListName "s")}}  {{.Short}}{{end}}{{end}}{{end}}
{{range .Plugins}}
    {{.Name | printf (print "%-" $.MaxRunListName "s")}}  {{.Short}} (plugin){{end}}
Run 'mc help [command]' for details.
Additional help topics:
{{range .Commands}}{{if not .Runnable}}
    {{.Name | printf "%-8s"}}  {{.Short}}{{end}}{{end}}
{{if .Dev}}This dev build of mc cannot auto-update itself.
{{end}}`[1:]))

var extraTemplate = template.Must(template.New("usage").Parse(`
Additional commands:
{{range .Commands}}{{if .Runnable}}{{if .ListAsExtra}}
    {{.Name | printf (print "%-" $.MaxRunExtraName "s")}}  {{.ShortExtra}}{{end}}{{end}}{{end}}
Run 'mc help [command]' for details.
`[1:]))

func printUsageTo(w io.Writer) {
	var runListNames []string
	for i := range commands {
		if commands[i].Runnable() && commands[i].List() {
			runListNames = append(runListNames, commands[i].Name())
		}
	}

	usageTemplate.Execute(w, struct {
		Commands       []*Command
		Dev            bool
		MaxRunListName int
	}{
		commands,
		Version == "dev",
		maxStrLen(runListNames),
	})
}
