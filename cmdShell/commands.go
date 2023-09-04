package cmdShell

import (
	"fmt"

	"github.com/MickMake/GoUnify/Only"
	"github.com/MickMake/GoUnify/cmdHelp"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)


const Group = "Run"

func (s *Shell) AttachCommands(cmd *cobra.Command) *cobra.Command {
	for range Only.Once {
		if cmd == nil {
			break
		}
		s.cmd = cmd

		name := aurora.White(s.name).Bold()	// c.cmd.Name()

		// ******************************************************************************** //
		s.SelfCmd = &cobra.Command{
			Use:                   "shell",
			Aliases:               []string{},
			Short:                 fmt.Sprintf("Run %s as an interactive shell.", name),
			Long:                  fmt.Sprintf("Run %s as an interactive shell.", name),
			DisableFlagParsing:    false,
			DisableFlagsInUseLine: false,
			PreRunE:               s.InitArgs,
			RunE:                  s.CmdHelpAll,
			Args:                  cobra.RangeArgs(0, 0),
		}
		cmd.AddCommand(s.SelfCmd)
		s.SelfCmd.Example = cmdHelp.PrintExamples(s.SelfCmd, "")
		s.SelfCmd.Annotations = map[string]string{"group": Group}

	}

	return s.SelfCmd
}

func (s *Shell) InitArgs(_ *cobra.Command, _ []string) error {
	for range Only.Once {
		//
	}
	return s.Error
}

func (s *Shell) CmdHelpAll(_ *cobra.Command, _ []string) error {
	for range Only.Once {
		s.Error = s.RunShell()
	}

	return s.Error
}
