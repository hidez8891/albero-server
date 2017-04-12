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

type zipFileReadCloser struct {
	r *zip.ReadCloser
	f io.ReadCloser
}

func (o *zipFileReadCloser) Read(p []byte) (n int, err error) {
	return o.f.Read(p)
}

func (o *zipFileReadCloser) Close() error {
	o.f.Close()
	return o.r.Close()
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
				r.Close()
				return nil
			}

			file, err := f.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				r.Close()
				return nil
			}

			zipfile := &zipFileReadCloser{
				r: r,
				f: file,
			}
			return conf.routing(zipfile, "", w, int64(f.UncompressedSize64))
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
