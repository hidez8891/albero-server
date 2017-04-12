package module

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"../utility/array"
)

type RoutingModule interface {
	ReturnFiles()
	ReturnBinary()
}

type moduleConfig struct {
	name    string
	exts    []string
	routing func(path, vpath string, w http.ResponseWriter) RoutingModule
}

var (
	confs []*moduleConfig
)

func SupportType() []string {
	exts := []string{}
	for _, conf := range confs {
		exts = append(exts, conf.exts...)
	}
	return exts
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

	// not found extension
	exti := strings.LastIndex(path, ".")
	if exti < 0 {
		http.Error(w, "Not Found File Type", http.StatusInternalServerError)
		return nil
	}
	ext := path[exti+1:]

	// search extension
	for _, conf := range confs {
		if array.IsInclude(ext, conf.exts) {
			return conf.routing(path, vpath, w)
		}
	}

	// not found file type
	http.Error(w, "No Support Type", http.StatusUnsupportedMediaType)
	return nil
}

func install(h *moduleConfig) *moduleConfig {
	confs = append(confs, h)
	return h
}

func returnBinary(w http.ResponseWriter, r io.Reader, mime string, size int64) {
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	if _, err := io.Copy(w, r); err != nil {
		log.Printf("ERR: %s: %v\n", mime, err)
	}
}
