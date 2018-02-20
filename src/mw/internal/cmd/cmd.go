/*
	Package cmd represents a command object.
*/
package cmd

import (
	"flag"
	"fmt"
)

// Commands holds the list of all the supported commands.
// This list is filled in main.go by referencing to it.
var Commands []*Command

var Usage func()

type Command struct {
	CmdName    string
	CmdUsage   string
	CmdFlag    flag.FlagSet
	HasNoFlags bool
	isGreat    bool
	Run        func(c *Command, args []string)
}

func (c *Command) Name() string {
	return c.CmdName
}

func (c *Command) Usage() {
	fmt.Println(c.CmdUsage)
}
