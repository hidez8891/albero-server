package plugin

import (
	"fmt"
	ospath "path"
)

// loaded plugin list
var plugins = map[string]Plugin{}

// Load is loaded plugin
func Load(path string) error {
	cmd, err := NewCmd(path)
	if err != nil {
		return err
	}

	pl := &_PluginBase{cmd: cmd}
	config, err := pl.Config()
	if err != nil {
		pl.close()
		return err
	}

	var plg Plugin
	switch config.Type {
	case TYPE_ARCH:
		plg = &PluginArch{pl}
	case TYPE_FILE:
		plg = &PluginFile{pl}
	}

	for _, ext := range config.Exts {
		plugins[ext] = plg
	}
	return nil
}

// Release is released plugins
func Release() {
	for _, pl := range plugins {
		pl.close()
	}
	plugins = map[string]Plugin{}
}

// Get returns plugin supported path file
func Get(path string) (Plugin, error) {
	ext := ospath.Ext(path)
	if len(ext) < 2 {
		return nil, fmt.Errorf("Path don't have file extension")
	}
	pl, ok := plugins[ext[1:]]
	if !ok {
		return nil, fmt.Errorf("No Support Type %s", ext)
	}
	return pl, nil
}
