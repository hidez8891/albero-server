package plugin

import (
	"fmt"

	"github.com/vmihailenco/msgpack"
)

// PluginArch :: Arch type plugin
type PluginArch struct {
	*_PluginBase
}

// Files return file list included path directory
func (o *PluginArch) Files(path string) ([]string, error) {
	block, err := msgpack.Marshal(path)
	if err != nil {
		return nil, err
	}
	block, err = o.call(FUNC_FILES, block)
	if err != nil {
		return nil, err
	}

	var files []string
	if err := msgpack.Unmarshal(block, &files); err != nil {
		return nil, err
	}
	return files, nil
}

// Files2 return file list included archive data
func (o *PluginArch) Files2(path string, data []byte) ([]string, error) {
	args := struct {
		path string
		data []byte
	}{path, data}
	block, err := msgpack.Marshal(&args)
	if err != nil {
		return nil, err
	}
	block, err = o.call(FUNC_FILES2, block)
	if err != nil {
		return nil, err
	}

	var files []string
	if err := msgpack.Unmarshal(block, &files); err != nil {
		return nil, err
	}
	return files, nil
}

// OpenFile return file's body and meta data from vpath
func (o *PluginArch) OpenFile(path, vpath string) (*File, error) {
	args := struct {
		path  string
		vpath string
	}{path, vpath}
	block, err := msgpack.Marshal(&args)
	if err != nil {
		return nil, err
	}
	block, err = o.call(FUNC_OPEN_FILE, block)
	if err != nil {
		return nil, err
	}

	var file File
	if err := msgpack.Unmarshal(block, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

// OpenFile2 return file's body and meta data form data
func (o *PluginArch) OpenFile2(path, vpath string, data []byte) (*File, error) {
	args := struct {
		path  string
		vpath string
		data  []byte
	}{path, vpath, data}
	block, err := msgpack.Marshal(&args)
	if err != nil {
		return nil, err
	}
	block, err = o.call(FUNC_OPEN_FILE2, block)
	if err != nil {
		return nil, err
	}

	var file File
	if err := msgpack.Unmarshal(block, &file); err != nil {
		return nil, err
	}
	return &file, nil
}

// Read :: No Supportted
func (o *PluginArch) Read(path string) (*File, error) {
	return nil, fmt.Errorf("Read has not Supportted this plugin")
}

// Read2 :: No Supportted
func (o *PluginArch) Read2(path string, data []byte) (*File, error) {
	return nil, fmt.Errorf("Read has not Supportted this plugin")
}
