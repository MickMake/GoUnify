package cmdDaemon

import (
	"github.com/MickMake/GoUnify/Only"
	"github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
)

type Daemon struct {
	cntxt *daemon.Context
	name  string

	Error error

	cmd     *cobra.Command
	SelfCmd *cobra.Command
}

func New(name string) *Daemon {
	var ret *Daemon

	for range Only.Once {
		ret = &Daemon {
			cntxt: &daemon.Context{},
			name: name,
			Error: nil,

			cmd:     nil,
			SelfCmd: nil,
		}
	}

	return ret
}

func (d *Daemon) GetCmd() *cobra.Command {
	return d.SelfCmd
}
