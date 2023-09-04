package cmdService

import (
	"github.com/MickMake/GoUnify/Only"
	"github.com/kardianos/service"
	"github.com/spf13/cobra"
)


type Service struct {
	name    string
	description string
	dir     string
	program   *program
	service   service.Service
	options   service.KeyValue
	config    service.Config

	Error     error

	cmd     *cobra.Command
	SelfCmd *cobra.Command
}

type program struct {
	exit chan struct{}
}


func New(name string, description string, configDir string) *Service {
	var ret *Service

	for range Only.Once {
		ret = &Service {
			name:    name,
			description: description,
			dir:     configDir,
			Error:   nil,

			cmd:     nil,
			SelfCmd: nil,
		}

		ret.options = make(service.KeyValue)
		// options["LaunchdConfig"] = fmt.Sprintf("%s%s%s%s",
		// 	"<key>StandardOutPath</key>\n",
		// 		"<string>/path/to/logfile.log</string>\n",
		// 		"<key>StandardErrorPath</key>\n",
		// 		"<string>/path/to/another_logfile.log</string>\n",
		// )
		ret.options["Restart"] = "on-success"
		ret.options["SuccessExitStatus"] = "1 2 8 SIGKILL"
		ret.options["KeepAlive"] = true
		ret.options["RunAtLoad"] = false
		ret.options["UserService"] = true
		ret.options["LogOutput"] = true
		// options["LogOutput"] = appConfig.Run.BaseDir + "/stdout.log"
		// options["StandardErrorPath"] = appConfig.Run.BaseDir + "/stderr.log"

		ret.config = service.Config {
			Name:        name,
			DisplayName: name,
			Description: description,
			// Dependencies: []string {
			// 	"Requires=network.target",
			// 	"After=network-online.target syslog.target",
			// },
			Option: ret.options,
			Arguments: []string{"run"},
			// Arguments: []string{"-daemon", "-config", ret.name + "/config.json"},
		}

		ret.program = &program{}
		// ret.service, ret.Error = service.New(ret.program, &ret.config)
		// if ret.Error != nil {
		// 	break
		// }
	}

	return ret
}

func (c *Service) GetCmd() *cobra.Command {
	return c.SelfCmd
}
