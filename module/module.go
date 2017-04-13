package module

import (
	"io"
	"log"
	"net/http"
	"os"
	ospath "path"
	"strconv"
	"strings"

	"../utility/array"
	"../utility/json"
)

type RoutingModule interface {
	ReturnFiles()
	ReturnBinary()
	Close()
}

type moduleConfig struct {
	name     string
	exts     []string
	types    moduleType
	routing  func(r io.ReadCloser, vpath string, w http.ResponseWriter, size int64) RoutingModule
	routing2 func(path, vpath string, w http.ResponseWriter) RoutingModule
}

type moduleType int

const (
	typeModuleImage moduleType = iota
	typeModuleArch
)

var (
	confs []*moduleConfig
)

func ReturnSupportType(w http.ResponseWriter, r *http.Request) {
	data := supportType()
	json.WriteResponse(w, data)
}

func Routing(path string, w http.ResponseWriter) RoutingModule {
	path = strings.Replace(path, "\\", "/", -1)

	// return directory routing module
	if stat, err := os.Stat(path); err == nil && stat.IsDir() {
		return newDirectoryRouting(path, w)
	}

	// split real path :: virtual path
	paths := strings.Split(path, "/")
	path, paths = paths[0], paths[1:]
	path += "/"

	for len(paths) > 0 {
		newPath := path + "/" + paths[0]
		if _, err := os.Stat(newPath); err != nil {
			break
		}

		path = newPath
		paths = paths[1:]
	}
	vpath := strings.Join(paths, "/")

	// invalid path
	if _, err := os.Stat(path); err != nil {
		http.Error(w, "Not Found File", http.StatusInternalServerError)
		return nil
	}

	// dispatch
	conf := dispatch(path, w)
	if conf == nil {
		return nil
	}
	return conf.routing2(path, vpath, w)
}

func install(h *moduleConfig) *moduleConfig {
	confs = append(confs, h)
	return h
}

func dispatch(path string, w http.ResponseWriter) *moduleConfig {
	// not found extension
	ext := ospath.Ext(path)
	if len(ext) == 0 {
		http.Error(w, "Not Found File Type", http.StatusInternalServerError)
		return nil
	}

	// search extension
	for _, conf := range confs {
		if array.IsInclude(ext, conf.exts) {
			return conf
		}
	}

	// not found file type
	http.Error(w, "No Support Type", http.StatusUnsupportedMediaType)
	return nil
}

func returnBinary(w http.ResponseWriter, r io.Reader, mime string, size int64) {
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	if _, err := io.Copy(w, r); err != nil {
		log.Printf("ERR: %s: %v\n", mime, err)
	}
}

func supportType() map[string][]string {
	files := []string{}
	archs := []string{}

	for _, conf := range confs {
		switch conf.types {
		case typeModuleImage:
			files = append(files, conf.exts...)
		case typeModuleArch:
			archs = append(archs, conf.exts...)
		default:
			// nothing to do
		}
	}

	return map[string][]string{
		"image":   files,
		"archive": archs,
	}
}
