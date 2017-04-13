package module

import (
	"io/ioutil"
	"net/http"
	"path"

	"../utility/array"
	"../utility/json"
)

type directoryRoutingModule struct {
	path string
	w    http.ResponseWriter
}

func (o *directoryRoutingModule) ReturnFiles() {
	paths, err := ioutil.ReadDir(o.path)
	if err != nil {
		http.Error(o.w, err.Error(), http.StatusInternalServerError)
		return
	}

	exts := supportType()
	dirs := []string{}
	archs := []string{}
	files := []string{}

	for _, pt := range paths {
		var (
			name = pt.Name()
			ext  = path.Ext(name)
		)

		switch {
		case pt.IsDir():
			dirs = append(dirs, name)
		case array.IsInclude(ext, exts["archive"]):
			archs = append(archs, name)
		case array.IsInclude(ext, exts["image"]):
			files = append(files, name)
		default:
			// nothing to do
		}
	}

	data := map[string][]string{
		"dir":     dirs,
		"image":   files,
		"archive": archs,
	}
	json.WriteResponse(o.w, data)
}

func (o *directoryRoutingModule) ReturnBinary() {
	http.Error(o.w, "Not Support", http.StatusUnsupportedMediaType)
}

func (o *directoryRoutingModule) Close() {
	http.Error(o.w, "Not Support", http.StatusUnsupportedMediaType)
}

func newDirectoryRouting(path string, w http.ResponseWriter) RoutingModule {
	return &directoryRoutingModule{
		path: path,
		w:    w,
	}
}
