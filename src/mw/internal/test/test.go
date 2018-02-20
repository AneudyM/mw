package test

import (
	"mw/internal/cmd"
)

var CmdTest = &cmd.Command{
	CmdName:    "test",
	CmdUsage:   "usage: mw test",
	HasNoFlags: true,
	Run:        testCode,
}

func testCode(c *cmd.Command, args []string) {

}
