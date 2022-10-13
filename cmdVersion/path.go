package cmdVersion

import (
	"bufio"
	"github.com/MickMake/GoUnify/Only"
	"github.com/MickMake/GoUnify/cmdLog"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)


func (v *Version) SetCmd(a ...string) error {
	var err error

	for range Only.Once {
		v.Cmd, err = filepath.Abs(filepath.Join(a...))
		if err != nil {
			break
		}

		v.CmdDir = filepath.Dir(v.Cmd)
		v.CmdFile = filepath.Base(v.Cmd)
	}

	return err
}

func (v *Version) IsBootstrapBinary() bool {
	var ok bool
	for range Only.Once {
		if v.ExecName != v.CmdFile {
			break
		}
		if v.ExecName != BootstrapBinaryName {
			break
		}
		ok = true
	}
	return ok
}

func (v *Version) AutoRun() State {
	for range Only.Once {
		if !v.AutoExec {
			break
		}

		if v.IsBootstrapBinary() {
			// Let's avoid an endless loop.
			break
		}

		if len(v.FullArgs) > 0 {
			if v.FullArgs[0] == CmdVersion {
				// Let's avoid another endless loop.
				break
			}
		}

		// @TODO - This is broken!
		cmd := exec.Command(v.TargetBinary, []string{"version", "info"}...)

		var stdout io.ReadCloser
		var err error
		stdout, err = cmd.StdoutPipe()
		if err != nil {
			break
		}

		err = cmd.Start()
		if err != nil {
			break
		}

		in := bufio.NewScanner(stdout)

		for in.Scan() {
			cmdLog.Printf(in.Text()) // write each line to your log, or anything you need
		}

		err = in.Err()
		if err != nil {
			cmdLog.Printf("error: %s", err)
			v.State.SetError(err.Error())
			break
		}

		// @TODO - This is broken!
		// fmt.Printf("Executing the real binary: '%s'\n", v.RuntimeBinary)
		// // c := exec.Command(v.TargetBinary, v.FullArgs...)
		// c := exec.Command(v.TargetBinary, []string{"version"}...)
		//
		// var stdoutBuf, stderrBuf bytes.Buffer
		// c.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
		// c.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
		// err := c.Run()
		// waitStatus := c.ProcessState.Sys().(syscall.WaitStatus)
		// waitStatus.ExitStatus()
		//
		// if err != nil {
		// 	fmt.Printf("stdoutBuf: %s\n", stdoutBuf.String())
		// 	fmt.Printf("stderrBuf: %s\n", stderrBuf.String())
		// 	v.State.SetError(err.Error())
		// 	break
		// }

		v.State.SetOk("")
	}

	return v.State
}

func (v *Version) CreateDummyBinary() State {
	for range Only.Once {
		var err error

		result := FileStat(v.RuntimeBinary, v.TargetBinary)
		if result.CopyOfRuntime {
			v.AutoExec = true
			break
		}

		//if result.IsRuntimeBinary {
		//	// We are running as the bootstrap binary.
		//	su.State.SetOk()
		//	break
		//}

		if result.LinkToRuntime {
			err = os.Remove(v.TargetBinary)
			if err != nil {
				v.State.SetError(err.Error())
				break
			}
			result.IsMissing = true
			v.AutoExec = true
		}

		if result.IsMissing {
			err = CopyFile(v.RuntimeBinary, v.TargetBinary)
			if err != nil {
				v.State.SetError(err.Error())
				break
			}
			v.AutoExec = true
		}
	}

	return v.State
}

func (v *Version) IsRunningAs(run string) bool {
	var ok bool
	// If OK - running executable file matches the string 'run'.
	//ok, err := regexp.MatchString("^" + run, r.CmdFile)

	if v.IsWindows() {
		//fmt.Printf("DEBUG: WINDOWS!\n")
		ok = strings.HasPrefix(run, strings.TrimSuffix(v.CmdFile, ".exe"))
		//run = strings.TrimSuffix(run, ".exe")
	} else {
		ok = strings.HasPrefix(run, v.CmdFile)
	}
	//fmt.Printf("DEBUG: Cmd.Runtime.IsRunningAs?? %s\n", ok)
	//fmt.Printf("DEBUG: run: %s\n", run)
	//fmt.Printf("DEBUG: r.CmdName: %s\n", r.CmdName)
	//fmt.Printf("DEBUG: r.CmdFile: %s\n", r.CmdFile)
	return ok
}

func (v *Version) IsRunningAsFile() bool {
	// If OK - running executable file matches the application binary name.
	//ok, err := regexp.MatchString("^" + r.CmdName, r.CmdFile)
	ok := strings.HasPrefix(v.ExecName, v.CmdFile)
	return ok
}

func (v *Version) IsRunningAsLink() bool {
	return !v.IsRunningAsFile()
}


type TargetFile struct {
	IsMissing       bool
	IsRuntimeBinary bool
	FileMatches     bool
	IsSymlink       bool
	LinkTo          string
	LinkEval        string
	LinkToRuntime   bool
	CopyOfRuntime   bool

	Error error
	Info  os.FileInfo
}

func FileStat(runtimeBinary string, targetBinary string) *TargetFile {
	var targetFile TargetFile

	for range Only.Once {
		targetFile.Info, targetFile.Error = os.Stat(targetBinary)
		if os.IsNotExist(targetFile.Error) {
			targetFile.IsMissing = true
		} else {
			targetFile.IsMissing = false

			if filepath.Base(runtimeBinary) == BootstrapBinaryName {
				targetFile.IsRuntimeBinary = true
			} else if runtimeBinary == targetBinary {
				targetFile.IsRuntimeBinary = true
				targetFile.CopyOfRuntime = true
			} else {
				targetFile.IsRuntimeBinary = false

				targetFile.Error = CompareBinary(runtimeBinary, targetBinary)
				if targetFile.Error == nil {
					targetFile.FileMatches = true
				} else {
					targetFile.FileMatches = false
				}
			}
		}

		targetFile.LinkTo, targetFile.Error = os.Readlink(targetBinary)
		if targetFile.LinkTo != "" {
			targetFile.IsSymlink = true

			targetFile.LinkEval, targetFile.Error = filepath.EvalSymlinks(targetBinary)
			if targetFile.LinkEval == "" {
				targetFile.LinkToRuntime = false
			} else {
				targetFile.LinkEval, targetFile.Error = filepath.Abs(targetFile.LinkEval)
				if targetFile.LinkEval == runtimeBinary {
					targetFile.LinkToRuntime = true
				} else if filepath.Base(targetFile.LinkEval) == BootstrapBinaryName {
					targetFile.LinkToRuntime = true
				} else {
					targetFile.LinkToRuntime = false
				}
			}
		} else {
			targetFile.IsSymlink = false
		}
	}

	return &targetFile
}
