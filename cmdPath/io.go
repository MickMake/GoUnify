package cmdPath

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MickMake/GoUnify/Only"
	"io/ioutil"
	"os"
)


// FileRead Retrieves data from a local file.
func (p *Path) FileRead(ref interface{}) error {
	var err error
	for range Only.Once {
		if p.name == "" {
			err = errors.New("empty file")
			break
		}

		var f *os.File
		f, err = os.Open(p.name)
		if err != nil {
			if os.IsNotExist(err) {
				err = nil
			}
			break
		}

		//goland:noinspection GoUnhandledErrorResult,GoDeferInLoop
		defer f.Close()

		err = json.NewDecoder(f).Decode(&ref)
	}

	// for range Only.Once {
	//	fn := ep.GetFilename()
	//	if err != nil {
	//		break
	//	}
	//
	//	ret, err = os.FileRead(fn)
	//	if err != nil {
	//		break
	//	}
	// }

	return err
}

// FileWrite Saves data to a file path.
func (p *Path) FileWrite(ref interface{}, perm os.FileMode) error {
	var err error
	for range Only.Once {
		if p.name == "" {
			err = errors.New("empty file")
			break
		}

		var f *os.File
		f, err = os.OpenFile(p.name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
		if err != nil {
			err = errors.New(fmt.Sprintf("Unable to write to file %s - %v", p.name, err))
			break
		}

		//goland:noinspection GoUnhandledErrorResult,GoDeferInLoop
		defer f.Close()
		err = json.NewEncoder(f).Encode(ref)

		// fn := ep.GetFilename()
		// if err != nil {
		//	break
		// }
		//
		// err = os.FileWrite(fn, data, perm)
		// if err != nil {
		//	break
		// }
	}

	return err
}

// PlainFileRead Retrieves data from a local file.
func (p *Path) PlainFileRead() ([]byte, error) {
	var data []byte
	var err error
	for range Only.Once {
		if p.name == "" {
			err = errors.New("empty file")
			break
		}

		var f *os.File
		f, err = os.Open(p.name)
		if err != nil {
			if os.IsNotExist(err) {
				err = nil
			}
			break
		}

		//goland:noinspection GoUnhandledErrorResult,GoDeferInLoop
		defer f.Close()

		data, err = ioutil.ReadAll(f)
	}

	return data, err
}

// PlainFileWrite Saves data to a file path.
func (p *Path) PlainFileWrite(data []byte, perm os.FileMode) error {
	var err error
	for range Only.Once {
		if p.name == "" {
			err = errors.New("empty file")
			break
		}

		var f *os.File
		f, err = os.OpenFile(p.name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
		if err != nil {
			err = errors.New(fmt.Sprintf("Unable to write to file %s - %v", p.name, err))
			break
		}
		//goland:noinspection GoUnhandledErrorResult,GoDeferInLoop
		defer f.Close()

		_, err = f.Write(data)
	}

	return err
}

// FileRemove Removes a file path.
func (p *Path) FileRemove() error {
	var err error
	for range Only.Once {
		if p.name == "" {
			err = errors.New("empty file")
			break
		}

		var f os.FileInfo
		f, err = os.Stat(p.name)
		if os.IsNotExist(err) {
			err = nil
			break
		}
		if err != nil {
			break
		}
		if f.IsDir() {
			err = errors.New("file is a directory")
			break
		}

		err = os.Remove(p.name)
	}

	return err
}
