package module

import (
	"archive/zip"
	"io"
	"net/http"
	ospath "path"
	"strings"

	"../utility/array"
	"../utility/json"
)

type archZipRoutingModule struct {
	path  string
	vpath string
	files []string
	w     http.ResponseWriter
}

func (o *archZipRoutingModule) ReturnFiles() {
	exts := supportType()
	dirs := []string{}
	files := []string{}
	archs := []string{}

	for _, f := range o.files {
		if strings.HasPrefix(f, o.vpath) {
			if ind := strings.Index(f, "/"); ind >= 0 {
				name := ospath.Join(o.path, f[:ind])
				dirs = append(dirs, name)
			} else {
				name := ospath.Join(o.path, f)
				ext := ospath.Ext(name)

				switch {
				case array.IsInclude(ext, exts["archive"]):
					// unsupport
				case array.IsInclude(ext, exts["image"]):
					files = append(files, name)
				default:
					// nothing to do
				}
			}
		}
	}

	data := map[string][]string{
		"dir":     array.Uniq(dirs),
		"image":   array.Uniq(files),
		"archive": array.Uniq(archs),
	}
	json.WriteResponse(o.w, data)
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
			path:  path,
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
	exts:     []string{".zip"},
	types:    typeModuleArch,
	routing:  archZipRouting,
	routing2: archZipRouting2,
})
