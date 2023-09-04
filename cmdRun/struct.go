package cmdRun

import (
	"github.com/MickMake/GoUnify/Only"
	"github.com/MickMake/GoUnify/cmdCron"
	"github.com/MickMake/GoUnify/cmdDaemon"
	"github.com/MickMake/GoUnify/cmdService"
	"github.com/MickMake/GoUnify/cmdShell"
	"github.com/spf13/cobra"
)


type Run struct {
	name        string
	version     string
	description string
	dir         string

	CmdDaemon  *cmdDaemon.Daemon
	CmdCron    *cmdCron.Cron
	CmdService *cmdService.Service
	CmdShell   *cmdShell.Shell

	Error       error

	cmd     *cobra.Command
	SelfCmd *cobra.Command
}

type program struct {
	exit chan struct{}
}


func New(name string, version string, description string, configDir string) *Run {
	var ret *Run

	for range Only.Once {
		ret = &Run {
			name:    name,
			version:    version,
			description: description,
			dir:     configDir,
			Error:   nil,

			cmd:     nil,
			SelfCmd: nil,
		}

		ret.CmdDaemon = cmdDaemon.New(name)
		ret.CmdCron = cmdCron.New(name)
		ret.CmdShell = cmdShell.New(name, version, configDir)
		ret.CmdService = cmdService.New(name, description, configDir)
	}

	return ret
}

func (c *Run) GetCmd() *cobra.Command {
	return c.SelfCmd
}
