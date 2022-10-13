package cmdPath

import (
	"bufio"
	"errors"
	"github.com/MickMake/GoUnify/Only"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)


type Path struct {
	name string
	err error
	fileinfo os.FileInfo
}

func NewPath(path ...string) *Path {
	var ret Path
	ret.Set(path...)
	return &ret
}

func (p *Path) DirExists() bool {
	var ok bool

	for range Only.Once {
		p.fileinfo, p.err = os.Stat(p.name)
		if os.IsNotExist(p.err) {
			break
		}

		if !p.fileinfo.IsDir() {
			break
		}

		ok = true
	}

	return ok
}

func (p *Path) FileExists() bool {
	var ok bool

	for range Only.Once {
		p.fileinfo, p.err = os.Stat(p.name)
		if os.IsNotExist(p.err) {
			break
		}

		if p.fileinfo.IsDir() {
			break
		}

		ok = true
	}

	return ok
}

func (p *Path) Chmod(mode os.FileMode) bool {
	var ok bool

	for range Only.Once {
		err := os.Chmod(p.name, mode)
		if err != nil {
			break
		}

		ok = true
	}

	return ok
}

func (p *Path) Set(path ...string) {

	for range Only.Once {
		dir := filepath.Join(path...)
		if strings.HasPrefix(dir, "~/") {
			u, err := user.Current()
			if err != nil {
				break
			}
			dir = strings.TrimPrefix(dir, "~/")
			dir = filepath.Join(u.HomeDir, dir)
		}

		p.name = dir
	}
}

func (p *Path) String() string {
	return p.name
}

func (p *Path) ModTime() time.Time {
	var ret time.Time

	for range Only.Once {
		if p.fileinfo == nil {
			p.fileinfo, p.err = os.Stat(p.name)
			if p.err != nil {
				break
			}
		}

		ret = p.fileinfo.ModTime()
	}
	return ret
}

//func (p *Path) Set(elem ...string) Path {
//	return (Path)(filepath.Join(elem...))
//}

func (p *Path) Join(elem ...string) Path {
	var pa []string
	//if p == nil {
	//	*p = "/"
	//}
	pa = append(pa, p.name)
	pa = append(pa, elem...)
	return Path{name: filepath.Join(pa...)}
}

func (p *Path) MkdirAll() error {
	var err error

	for range Only.Once {
		if p.DirExists() {
			break
		}

		dir := p.name
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			break
		}

		if !p.DirExists() {
			err = errors.New("no dir")
			break
		}
	}

	return err
}

func (p *Path) Copy(fp string) error {

	for range Only.Once {
		var stat os.FileInfo
		p.fileinfo, p.err = os.Stat(p.name)
		if os.IsNotExist(p.err) {
			break
		}
		if stat.IsDir() {
			p.err = errors.New("file is a dir")
			break
		}

		var input []byte
		input, p.err = ioutil.ReadFile(fp)
		if p.err != nil {
			break
		}

		dfp := filepath.Join(p.name, filepath.Base(fp))
		p.err = ioutil.WriteFile(dfp, input, stat.Mode())
		if p.err != nil {
			break
		}
	}

	return p.err
}

func (p *Path) Move(fp string) error {
	var err error

	for range Only.Once {
		err = p.Copy(fp)
		if err != nil {
			break
		}

		err = os.Remove(fp)
		if err != nil {
			break
		}
	}

	return err
}

func (p *Path) GrepFile(search string) (int, error) {
	var line int
	var err error

	for range Only.Once {
		p.Set(p.name)

		var f *os.File
		f, err = os.Open(p.name)
		if err != nil {
			// Silently ignore missing files.
			err = nil
			break
		}
		//goland:noinspection GoDeferInLoop,GoUnhandledErrorResult
		defer f.Close()

		// Splits on newlines by default.
		scanner := bufio.NewScanner(f)
		line = 1
		// https://golang.org/pkg/bufio/#Scanner.Scan
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), search) {
				break
			}

			line++
		}

		err = scanner.Err()
		if err != nil {
			break
		}
	}

	return line, err
}


//goland:noinspection SpellCheckingInspection
var RcFiles = []Path{
	// BASH
	{name: "/etc/profile"},
	{name: "/etc/bashrc"},
	{name: "~/.profile"},
	{name: "~/.bash_profile"},
	{name: "~/.bashrc"},
	{name: "~/.bash_login"},
	{name: "~/.bash_logout"},

	// ZSH
	{name: "/etc/zlogin"},
	{name: "/etc/zlogout"},
	{name: "/etc/zprofile"},
	{name: "/etc/zshenv"},
	{name: "/etc/zshrc"},
	{name: "~/.zlogin"},
	{name: "~/.zlogout"},
	{name: "~/.zprofile"},
	{name: "~/.zshenv"},
	{name: "~/.zshrc"},

	// CSH
	{name: "/etc/csh.cshrc"},
	{name: "/etc/csh.login"},
	{name: "/etc/csh.logout"},
	{name: "~/.cshrc"},
	{name: "~/.login"},
	{name: "~/.logout"},
}

//goland:noinspection GoUnusedExportedFunction
func GrepFiles(search string, fps ...Path) ([]string, error) {
	var files []string
	var err error

	if fps == nil {
		fps = RcFiles
	}
	if len(fps) == 0 {
		fps = RcFiles
	}

	for _, p := range fps {
		var line int
		line, err = p.GrepFile(search)
		if line > 0 {
			files = append(files, p.String()+" line:"+strconv.Itoa(line))
		}
	}

	return files, err
}

func ResolveFile(file string) string {
	var result string
	var err error

	for range Only.Once {
		_, err = os.Stat(file)
		if os.IsNotExist(err) {
			break
		}

		result, err = os.Readlink(file)
		if result == "" {
			result = file
			break
		}

		result, err = filepath.EvalSymlinks(file)
		if result == "" {
			result = file
			break
		}

		result, err = filepath.Abs(result)
		if result == "" {
			result = file
			break
		}
	}

	return result
}
