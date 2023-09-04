package cmdShell

import (
	"fmt"
	"github.com/MickMake/GoUnify/Only"
	"github.com/MickMake/GoUnify/cmdExec"
	"github.com/abiosoft/ishell/v2"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"path/filepath"
	"time"
)


type Shell struct {
	name    string
	version string
	history string
	*ishell.Shell

	cmd     *cobra.Command
	SelfCmd *cobra.Command
	Error error
}

func New(name string, version string, configDir string) *Shell {
	var ret Shell
	for range Only.Once {
		ret = Shell {
			name: name,
			version: version,
			history: filepath.Join(configDir, "history"),
			Shell: ishell.New(),
			Error: nil,
		}

		// shell.SetHomeHistoryPath(".ishell_history")
		ret.SetHistoryPath(ret.history)
		ret.AutoHelp(true)
		ret.IgnoreCase(false)
		ret.SetPager("less", []string{"-SinR"})
		ret.SetPrompt(fmt.Sprintf("%s v%s # ", ret.name, ret.version))
		// ret.ShowPaged("")
		ret.ShowPrompt(true)
	}
	return &ret
}

func (s *Shell) RunShell() error {
	for range Only.Once {
		s.cmd.Hidden = true		// Stop the command parent from being run within a shell.
		s.BuildCmd(s.cmd.Root(), nil)
		s.AddCmd(&ishell.Cmd {
			Name:     "spinner",
			Aliases:  []string{},
			Func: func(c *ishell.Context) {
				display := ishell.ProgressDisplayCharSet(spinner.CharSets[38])
				fmt.Printf("%v\n", display)
				c.ProgressBar().Display(display)
				c.ProgressBar().Start()
				c.ProgressBar().Final("|")
				c.Printf("|")
				for i := 0; i < 101; i++ {
					c.ProgressBar().Prefix(fmt.Sprint("-", i, "%"))
					c.ProgressBar().Progress(i)
					time.Sleep(time.Millisecond * 100) // some background computation
				}
				c.ProgressBar().Stop()
			},
			Help:     "",
			LongHelp: "",
		})

		s.AddCmd(&ishell.Cmd {
			Name:     "cls",
			Aliases:  []string{},
			Func: func(c *ishell.Context) {
				c.ClearScreen()
			},
			Help:     "",
			LongHelp: "",
		})

		s.AddCmd(&ishell.Cmd {
			Name:     "?",
			Aliases:  []string{},
			Func: func(c *ishell.Context) {
				s.cmd.Hidden = true		// Stop the command parent from being run within a shell.
				_ = s.cmd.Root().Help()
			},
			Help:     "",
			LongHelp: "",
		})

		s.Interrupt(func(c *ishell.Context, count int, input string) {
			// fmt.Printf("[%d]:[%s]: %v\n", count, input, c)
			switch {
				case count == 1:
					fmt.Println("Ctrl-c once more to exit")
				case count > 1:
					c.Stop()
					// os.Exit(1)
			}
		})


		s.Run()
		fmt.Println("Terminated")
		s.Wait()
	}
	return s.Error
}

func (s *Shell) CmdFunc(c *ishell.Context) {
	for range Only.Once {
		s.Error = s.ReparseArgs(c.RawArgs...)
		if s.Error != nil {
			break
		}
	}
}

func (s *Shell) ReparseArgs(args ...string) error {
	for range Only.Once {
		cmdExec.ResetArgs(args...)

		// rootCmd := cmdExec.FindRoot(cmd)
		// rootCmd.SetArgs(os.Args)
		s.cmd.Hidden = true		// Stop the command parent from being run within a shell.
		s.Error = s.cmd.Root().Execute()
		if s.Error != nil {
			break
		}
	}

	return s.Error
}

func (s *Shell) BuildCmd(cmd *cobra.Command, parent *ShellCmd) *ShellCmd {
	for range Only.Once {
		if parent == nil {
			parent = &ShellCmd{
				Ishell: &ishell.Cmd {
					Name:     cmd.Name(),
					Aliases:  cmd.Aliases,
					Func:     s.CmdFunc,
					Help:     cmd.CommandPath() + " " + cmd.Short,
					LongHelp: cmd.UseLine(),
				},
				Cobra: cmd,
			}
			s.SetRootCmd(parent.Ishell)
			s.Println(cmd.Long)
		}

		for _, c := range cmd.Commands() {
			// if c.Name() == "help" {
			// 	continue
			// }
			child := &ShellCmd{
				Ishell: &ishell.Cmd {
					Name:     c.Name(),
					Aliases:  c.Aliases,
					Func:     s.CmdFunc,
					Help:     c.CommandPath() + " " + c.Short,
					LongHelp: c.UseLine(),
				},
				Cobra: c,
			}
			if len(c.Commands()) == 0 {
				// fmt.Printf("%s %s - %s - %s\n", retString, c.Name(), c.UseLine(), c.Use)
				parent.Ishell.AddCmd(child.Ishell)
				continue
			}

			child = s.BuildCmd(c, child)
			parent.Ishell.AddCmd(child.Ishell)
		}
	}
	return parent
}


type ShellCmd struct {
	Ishell *ishell.Cmd
	Cobra *cobra.Command
}
