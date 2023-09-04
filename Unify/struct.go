// Package Unify - This package contains common functionality that's used across multiple binaries.
// It's an easy way to include some important functionality into every binary.
// Currently, it provides:
// - Cron scheduler.
// - Daemonizing a process.
// - Logging.
// - Version control and self-update.
// - Cobra/Viper integration.
// - Interactive shell based on Cobra CLI commands.
package Unify

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/MickMake/GoUnify/Only"
	"github.com/MickMake/GoUnify/cmdConfig"
	"github.com/MickMake/GoUnify/cmdCron"
	"github.com/MickMake/GoUnify/cmdDaemon"
	"github.com/MickMake/GoUnify/cmdHelp"
	"github.com/MickMake/GoUnify/cmdRun"
	"github.com/MickMake/GoUnify/cmdService"
	"github.com/MickMake/GoUnify/cmdShell"
	"github.com/MickMake/GoUnify/cmdVersion"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


type Unify struct {
	Options  Options  `json:"options"`
	Flags    Flags    `json:"flags"`
	Commands Commands `json:"commands"`

	mergeRun bool
	Error error `json:"-"`
}

type Options struct {
	Description   string `json:"description"`
	BinaryName    string `json:"binary_name"`
	BinaryVersion string `json:"binary_version"`
	SourceRepo    string `json:"source_repo"`
	BinaryRepo    string `json:"binary_repo"`
	EnvPrefix     string `json:"env_prefix"`
	HelpSummary   string `json:"help_template"`
	ReadMe        string `json:"readme"`
	Examples      string `json:"examples"`
}

type Flags struct {
	ConfigFile string        `json:"config_file"`
	ConfigDir  string        `json:"config_dir"`
	CacheDir   string        `json:"cache_dir"`
	Quiet      bool          `json:"quiet"`
	Debug      bool          `json:"debug"`
	Timeout    time.Duration `json:"timeout"`
	MergeRun   bool          `json:"merge_run"`
}


// New - Create new Unify instance.
func New(options Options, flags Flags) *Unify {
	var unify Unify

	for range Only.Once {
		unify.Options = options
		unify.Flags = flags

		if unify.Options.EnvPrefix == "" {
			unify.Options.EnvPrefix = "UNIFY"	// cmdVersion.GetEnvPrefix()
		}
		unify.Options.EnvPrefix = unify.GetEnvPrefix()

		// MergeRun == true - Merge the commands "daemon", "cron", "shell", "service"
		unify.Error = unify.InitCmds(flags.MergeRun)
		if unify.Error != nil {
			break
		}

		unify.Error = unify.InitFlags()
		if unify.Error != nil {
			break
		}
	}

	return &unify
}

// InitCmds -
func (u *Unify) InitCmds(mergeRun bool) error {

	for range Only.Once {
		// ******************************************************************************** //
		u.Commands.CmdRoot = &cobra.Command {
			Use:              u.Options.BinaryName,
			Short:            fmt.Sprintf("%s - %s", u.Options.BinaryName, u.Options.Description),
			Long:             fmt.Sprintf("%s - %s", u.Options.BinaryName, u.Options.Description),
			RunE:             u.CmdRoot,
			TraverseChildren: true,
		}
		u.Commands.CmdRoot.Example = cmdHelp.PrintExamples(u.Commands.CmdRoot, "")

		u.Commands.CmdConfig = cmdConfig.New(u.Options.BinaryName, u.Options.BinaryVersion, u.Options.EnvPrefix)

		u.Commands.CmdVersion = cmdVersion.New(u.Options.BinaryName, u.Options.BinaryVersion, false)
		u.Commands.CmdVersion.SetBinaryRepo(u.Options.BinaryRepo)
		u.Commands.CmdVersion.SetSourceRepo(u.Options.SourceRepo)

		u.mergeRun = mergeRun
		if u.mergeRun {
			u.Commands.CmdRun = cmdRun.New(u.Options.BinaryName, u.Options.BinaryVersion, u.Options.Description, u.GetConfigDir())
		} else {
			u.Commands.CmdDaemon = cmdDaemon.New(u.Options.BinaryName)
			u.Commands.CmdCron = cmdCron.New(u.Options.BinaryName)
			u.Commands.CmdShell = cmdShell.New(u.Options.BinaryName, u.Options.BinaryVersion, u.GetConfigDir())
			u.Commands.CmdService = cmdService.New(u.Options.BinaryName, u.Options.Description, u.GetConfigDir())
			// u.Commands.CmdSystray = cmdSystray.New(u.Commands.CmdConfig, u.Commands.CmdVersion)
		}

		u.Commands.CmdHelp = cmdHelp.New()
		u.Commands.CmdHelp.SetCommand(u.Options.BinaryName)
		u.Commands.CmdHelp.SetHelpSummary(u.Options.HelpSummary)
		u.Commands.CmdHelp.SetEnvPrefix(u.Options.EnvPrefix)
		u.Commands.CmdHelp.SetReadMe(u.Options.ReadMe)
		u.Commands.CmdHelp.SetExamples(u.Options.Examples)
	}

	return u.Error
}

// InitFlags -
func (u *Unify) InitFlags() error {

	for range Only.Once {
		u.Commands.CmdRoot.PersistentFlags().StringVar(&u.Flags.ConfigFile, cmdConfig.ConfigFileFlag, defaultConfig, fmt.Sprintf("%s: config file.", u.Options.BinaryName))
		// _ = rootCmd.PersistentFlags().MarkHidden(flagConfigFile)
		u.Commands.CmdRoot.PersistentFlags().BoolVarP(&u.Flags.Debug, flagDebug, "", defaultDebug, fmt.Sprintf("%s: Debug mode.", u.Options.BinaryName))
		u.Commands.CmdConfig.SetDefault(flagDebug, false)
		u.Commands.CmdRoot.PersistentFlags().BoolVarP(&u.Flags.Quiet, flagQuiet, "", defaultQuiet, fmt.Sprintf("%s: Silence all messages.", u.Options.BinaryName))
		u.Commands.CmdConfig.SetDefault(flagQuiet, false)
		u.Commands.CmdRoot.PersistentFlags().DurationVarP(&u.Flags.Timeout, flagTimeout, "", defaultTimeout, fmt.Sprintf("Web timeout."))
		u.Commands.CmdConfig.SetDefault(flagTimeout, defaultTimeout)

		u.Commands.CmdRoot.PersistentFlags().SortFlags = false
		u.Commands.CmdRoot.Flags().SortFlags = false

		// cobra.OnInitialize(initConfig)	// Bound to rootCmd now.
		cobra.EnableCommandSorting = false
	}

	return u.Error
}

// Execute -
func (u *Unify) Execute() error {

	for range Only.Once {
		if u.mergeRun {
			u.Commands.CmdRun.AttachCommands(u.Commands.CmdRoot)
		} else {
			u.Commands.CmdDaemon.AttachCommands(u.Commands.CmdRoot)
			u.Commands.CmdCron.AttachCommands(u.Commands.CmdRoot)
			u.Commands.CmdShell.AttachCommands(u.Commands.CmdRoot)
			u.Commands.CmdService.AttachCommands(u.Commands.CmdRoot)
			// u.Commands.CmdSystray.AttachCommands(u.Commands.CmdRoot)
		}

		u.Commands.CmdConfig.AttachCommands(u.Commands.CmdRoot)
		u.Commands.CmdVersion.AttachCommands(u.Commands.CmdRoot, true)
		u.Commands.CmdHelp.AttachCommands(u.Commands.CmdRoot)

		if u.Flags.ConfigDir != "" {
			u.Commands.CmdConfig.SetDir(u.Flags.ConfigDir)
		} else {
			u.Flags.ConfigDir = u.Commands.CmdConfig.Dir
		}

		if u.Flags.ConfigFile != "" {
			u.Commands.CmdConfig.SetFile(u.Flags.ConfigFile)
		} else {
			u.Flags.ConfigFile = u.Commands.CmdConfig.File
		}

		if u.Flags.CacheDir == "" {
			u.Flags.CacheDir = u.GetCacheDir()
		}

		if u.Flags.Timeout == 0 {
			u.Flags.Timeout = defaultTimeout
		}

		u.Commands.CmdRoot.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return u.Commands.CmdConfig.Init(cmd)
		}

		cc.Init(&cc.Config{
			RootCmd:         u.Commands.CmdRoot,
			Headings:        cc.HiCyan + cc.Bold + cc.Underline,
			Commands:        cc.HiYellow + cc.Bold,
			Example:         cc.Italic,
			ExecName:        cc.Bold,
			Flags:           cc.Bold,
			NoExtraNewlines: true,
			NoBottomNewline: true,
			CmdShortDescr:   0,
			FlagsDataType:   0,
			FlagsDescr:      0,
			Aliases:         0,
		})

		u.Error = u.Commands.Execute()
		if u.Error != nil {
			break
		}
	}
	return u.Error
}

// GetCmd -
func (u *Unify) GetCmd() *cobra.Command {
	return u.Commands.CmdRoot
}

// GetViper -
func (u *Unify) GetViper() *viper.Viper {
	return u.Commands.CmdConfig.GetViper()
}

// WriteConfig -
func (u *Unify) WriteConfig() error {
	return u.Commands.CmdConfig.Write()
}

// ReadConfig -
func (u *Unify) ReadConfig() error {
	return u.Commands.CmdConfig.Read()
}

// GetConfigDir -
func (u *Unify) GetConfigDir() string {
	return u.Commands.CmdConfig.Dir
}

// GetConfigFile -
func (u *Unify) GetConfigFile() string {
	return u.Commands.CmdConfig.File
}

// GetCacheDir -
func (u *Unify) GetCacheDir() string {
	return filepath.Join(u.Commands.CmdConfig.Dir, "cache")
}


type Commands struct {
	CmdRoot    *cobra.Command
	CmdRun     *cmdRun.Run

	CmdDaemon  *cmdDaemon.Daemon
	CmdCron    *cmdCron.Cron
	CmdService *cmdService.Service
	CmdShell   *cmdShell.Shell

	CmdVersion *cmdVersion.Version
	CmdConfig  *cmdConfig.Config
	// CmdSystray *cmdSystray.Config
	CmdHelp    *cmdHelp.Help
}

// Execute -
func (c *Commands) Execute() error {
	return c.CmdRoot.Execute()
}


// CmdRoot -
func (u *Unify) CmdRoot(cmd *cobra.Command, args []string) error {
	for range Only.Once {
		if len(args) == 0 {
			_ = cmd.Help()
			break
		}
		u.Error = errors.New(fmt.Sprintf("Unknown command string: %v\n", args))
	}
	return u.Error
}
