package main

import (
	"fmt"
	"os"
	"strings"

	flag "github.com/bgentry/pflag"
)

type Command struct {
	Run         func(cmd *Command, args []string)
	Flag        flag.FlagSet
	NeedsServer bool

	Usage    string // first word must be command name
	Category string
	Short    string
	Long     string
}

func (c *Command) PrintUsage() {
	if c.Runnable() {
		fmt.Fprintf(os.Stderr, "Usage: mc %s\n", c.FullUsage())
	}
	fmt.Fprintf(os.Stderr, "Use 'mc help %s' for more information.\n", c.Name())
}

func (c *Command) PrintLongUsage() {
	if c.Runnable() {
		fmt.Fprintf(os.Stderr, "Usage: mc %s\n", c.FullUsage())
	}
	fmt.Println(strings.Trim(c.Long, "\n"))
}

func (c *Command) FullUsage() string {
	return c.Usage
}

func (c *Command) Name() string {
	name := c.Usage
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}

const extra = " (extra)"

func (c *Command) List() bool {
	return c.Short != "" && !strings.HasSuffix(c.Short, extra)
}

func (c *Command) ListAsExtra() bool {
	return c.Short != "" && strings.HasSuffix(c.Short, extra)
}

func (c *Command) ShortExtra() string {
	return c.Short[:len(c.Short)-len(extra)]
}
