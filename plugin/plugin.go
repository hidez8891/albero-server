package plugin

import (
	"github.com/vmihailenco/msgpack"
)

// function ID
type FuncID int

const (
	FUNC_CONFIG FuncID = iota
	FUNC_CLOSE
	FUNC_FILES
	FUNC_FILES2
	FUNC_OPEN_FILE
	FUNC_OPEN_FILE2
	FUNC_READ
	FUNC_READ2
)

// Plugin Type
type PluginType int

const (
	TYPE_ARCH PluginType = iota
	TYPE_FILE
)

// Plugin config
type Config struct {
	Type PluginType
	Exts []string
}

// File struct
type File struct {
	Size int64
	Mime string
	Data []byte
}

//
// Plugin
//
type Plugin interface {
	Config() (*Config, error)
	close() error

	// archive type
	Files(string) ([]string, error)
	Files2(string, []byte) ([]string, error)
	OpenFile(string, string) (*File, error)
	OpenFile2(string, string, []byte) (*File, error)

	// file type
	Read(string) (*File, error)
	Read2(string, []byte) (*File, error)
}

//
// Based Plugin
//
type _PluginBase struct {
	cmd *Cmd
}

// Config return plugin config
func (o *_PluginBase) Config() (*Config, error) {
	block, err := o.call(FUNC_CONFIG, []byte{})
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := msgpack.Unmarshal(block, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

// close is killed plugin process
func (o *_PluginBase) close() error {
	if err := o.cmd.SendID(FUNC_CLOSE); err != nil {
		return err
	}
	if err := o.cmd.SendArgs([]byte{}); err != nil {
		return err
	}
	return nil
}

// call is called plugin process function
func (o *_PluginBase) call(id FuncID, args []byte) ([]byte, error) {
	if err := o.cmd.SendID(id); err != nil {
		return nil, err
	}
	if err := o.cmd.SendArgs(args); err != nil {
		return nil, err
	}
	return o.cmd.RecvReturn()
}
