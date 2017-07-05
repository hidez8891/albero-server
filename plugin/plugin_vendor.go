package plugin

import (
	"fmt"
	"os"

	"github.com/vmihailenco/msgpack"
)

// DispatchArchLoop dispatch functions for arch type
func DispatchArchLoop(adapt PluginArchAdapter) error {
	cmd, err := NewCmd2(os.Stdin, os.Stdout)
	if err != nil {
		return err
	}

	for {
		funcID, err := cmd.RecvID()
		if err != nil {
			return err
		}
		args, err := cmd.RecvArgs()
		if err != nil {
			return err
		}

		switch funcID {
		case FUNC_CONFIG:
			conf, err := adapt.Config()
			if err != nil {
				return err
			}
			rets, err := msgpack.Marshal(conf)
			if err != nil {
				return err
			}
			if err := cmd.SendReturn(rets); err != nil {
				return err
			}

		case FUNC_CLOSE:
			return nil

		case FUNC_FILES:
			var path string
			if err := msgpack.Unmarshal(args, &path); err != nil {
				return err
			}
			files, err := adapt.Files(path)
			if err != nil {
				return err
			}
			rets, err := msgpack.Marshal(files)
			if err != nil {
				return err
			}
			if err := cmd.SendReturn(rets); err != nil {
				return err
			}

		case FUNC_FILES2:
			fargs := struct {
				path string
				data []byte
			}{}
			if err := msgpack.Unmarshal(args, &fargs); err != nil {
				return err
			}
			files, err := adapt.Files2(fargs.path, fargs.data)
			if err != nil {
				return err
			}
			rets, err := msgpack.Marshal(files)
			if err != nil {
				return err
			}
			if err := cmd.SendReturn(rets); err != nil {
				return err
			}

		case FUNC_OPEN_FILE:
			fargs := struct {
				path  string
				vpath string
			}{}
			if err := msgpack.Unmarshal(args, &fargs); err != nil {
				return err
			}
			file, err := adapt.OpenFile(fargs.path, fargs.vpath)
			if err != nil {
				return err
			}
			rets, err := msgpack.Marshal(file)
			if err != nil {
				return err
			}
			if err := cmd.SendReturn(rets); err != nil {
				return err
			}

		case FUNC_OPEN_FILE2:
			fargs := struct {
				path  string
				vpath string
				data  []byte
			}{}
			if err := msgpack.Unmarshal(args, &fargs); err != nil {
				return err
			}
			file, err := adapt.OpenFile2(fargs.path, fargs.vpath, fargs.data)
			if err != nil {
				return err
			}
			rets, err := msgpack.Marshal(file)
			if err != nil {
				return err
			}
			if err := cmd.SendReturn(rets); err != nil {
				return err
			}

		default:
			return fmt.Errorf("called Un-Support funcID %d", funcID)
		}
	}
}

// Arch type adapter
type PluginArchAdapter interface {
	Config() (*Config, error)
	close() error
	Files(string) ([]string, error)
	Files2(string, []byte) ([]string, error)
	OpenFile(string, string) (*File, error)
	OpenFile2(string, string, []byte) (*File, error)
}

// File type adapter
type PluginFileAdapter interface {
	Config() (*Config, error)
	close() error
	Read(string) (*File, error)
	Read2(string, []byte) (*File, error)
}
