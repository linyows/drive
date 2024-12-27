package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/linyows/probe"
	"github.com/linyows/probe/actions/hello"
	http "github.com/linyows/probe/actions/http"
	"github.com/linyows/probe/actions/smtp"
)

type Cmd struct {
	WorkflowPath string
	Init         bool
	Lint         bool
	Help         bool
	Verbose      bool
	validFlags   []string
	ver          string
	rev          string
}

func runBuiltinActions(name string) {
	switch name {
	case "http":
		http.Serve()
	case "hello":
		hello.Serve()
	case "smtp":
		smtp.Serve()
	}
}

func newCmd(args []string) *Cmd {
	if len(args) >= 3 && args[1] == probe.BuiltinCmd {
		runBuiltinActions(args[2])
		return nil
	}

	c := Cmd{
		validFlags: []string{"help", "init", "lint", "workflow", "verbose"},
		ver:        version,
		rev:        commit,
	}

	flag.StringVar(&c.WorkflowPath, "workflow", "", "Specify yaml-path of workflow")
	flag.BoolVar(&c.Help, "help", false, "Show command usage")
	flag.BoolVar(&c.Init, "init", false, "Export a workflow template as yaml file")
	flag.BoolVar(&c.Lint, "lint", false, "Check the syntax in workflow")
	flag.BoolVar(&c.Verbose, "verbose", false, "Show verbose log")

	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "-") && !c.isValid(arg) {
			fmt.Printf("Unknown flag: %s\n", arg)
			fmt.Println("try --help to know more")
			return nil
		}
	}

	flag.Parse()
	return &c
}

func (c *Cmd) isValid(flag string) bool {
	if idx := strings.Index(flag, "="); idx != -1 {
		flag = flag[:idx]
	}

	for _, validFlag := range c.validFlags {
		if strings.TrimLeft(flag, "-") == validFlag {
			return true
		}
	}

	return false
}

func (c *Cmd) usage() {
	h := `
Probe - scenario testing tool (ver: %s [%s])

Usage: probe [options] <command>
`
	h = strings.TrimPrefix(h, "\n")
	fmt.Fprint(flag.CommandLine.Output(), fmt.Sprintf(h, c.ver, c.rev))
}

func (c *Cmd) start() {
	switch {
	case c.Help:
		c.usage()
	case c.Lint:
	case c.Init:
	default:
		p := probe.New(c.WorkflowPath, c.Verbose)
		if err := p.Do(); err != nil {
			fmt.Printf("%#v\n", err)
		}
	}
}
