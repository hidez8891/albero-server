package plugin

import (
	"fmt"

	"github.com/vmihailenco/msgpack"
)

// PluginFile :: File type plugin
type PluginFile struct {
	*_PluginBase
}

// Files :: No Supportted
func (o *PluginFile) Files(path string) ([]string, error) {
	return nil, fmt.Errorf("Read has not Supportted this plugin")
}

// Files2 :: No Supportted
func (o *PluginFile) Files2(path string, data []byte) ([]string, error) {
	return nil, fmt.Errorf("Read has not Supportted this plugin")
}

// OpenFile :: No Supportted
func (o *PluginFile) OpenFile(path, vpath string) (*File, error) {
	return nil, fmt.Errorf("Read has not Supportted this plugin")
}

// OpenFile2 :: No Supportted
func (o *PluginFile) OpenFile2(path, vpath string, data []byte) (*File, error) {
	return nil, fmt.Errorf("Read has not Supportted this plugin")
}

// Read returns file data
func (o *PluginFile) Read(path string) (*File, error) {
	block, err := msgpack.Marshal(&path)
	if err != nil {
		return nil, err
	}
	block, err = o.call(FUNC_READ, block)
	if err != nil {
		return nil, err
	}

	var file File
	if err := msgpack.Unmarshal(block, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

// Read2 returns file data
func (o *PluginFile) Read2(path string, data []byte) (*File, error) {
	args := struct {
		path string
		data []byte
	}{path, data}
	block, err := msgpack.Marshal(&args)
	if err != nil {
		return nil, err
	}
	block, err = o.call(FUNC_READ2, block)
	if err != nil {
		return nil, err
	}

	var file File
	if err := msgpack.Unmarshal(block, &file); err != nil {
		return nil, err
	}
	return &file, nil
}
