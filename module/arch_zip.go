package module

import (
	"archive/zip"
	"io"
	"net/http"
	"strings"

	"../utility/array"
)

type archZipRoutingModule struct {
	vpath string
	files []string
	w     http.ResponseWriter
}

func (o *archZipRoutingModule) ReturnFiles() {
	// TODO
}

func (o *archZipRoutingModule) ReturnBinary() {
	http.Error(o.w, "Not Support", http.StatusUnsupportedMediaType)
}

func (o *archZipRoutingModule) Close() {
	// Nothing to do
}

func archZipRouting(r io.ReadCloser, vpath string, w http.ResponseWriter, size int64) RoutingModule {
	http.Error(w, "Not Support", http.StatusUnsupportedMediaType)
	r.Close()
	return nil
}

func archZipRouting2(path, vpath string, w http.ResponseWriter) RoutingModule {
	// open archive
	r, err := zip.OpenReader(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	defer r.Close()

	// get all file path
	paths := make([]string, len(r.File))
	for i, f := range r.File {
		paths[i] = strings.Replace(f.Name, "\\", "/", -1)
	}

	// vpath is include file ?
	if vpath != "" {
		for _, f := range r.File {
			if f.Name != vpath {
				continue
			}

			conf := dispatch(vpath, w)
			if conf == nil {
				return nil
			}

			file, err := f.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil
			}
			return conf.routing(file, "", w, int64(f.UncompressedSize64))
		}
	}

	// vpath is include directory ?
	if vpath == "" || array.IsIncludeFunc(vpath, paths, strings.HasPrefix) {
		return &archZipRoutingModule{
			vpath: vpath,
			files: paths,
			w:     w,
		}
	}

	http.Error(w, "Not Found", http.StatusInternalServerError)
	return nil
}

var zipConf = install(&moduleConfig{
	name:     "zip",
	exts:     []string{"zip"},
	routing:  archZipRouting,
	routing2: archZipRouting2,
})
