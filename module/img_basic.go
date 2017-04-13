package module

import (
	"io"
	"net/http"
	"os"
)

type imgBasicRoutingModule struct {
	mime string
	w    http.ResponseWriter
	r    io.ReadCloser
	size int64
}

func (o *imgBasicRoutingModule) ReturnFiles() {
	http.Error(o.w, "Not Support", http.StatusUnsupportedMediaType)
}

func (o *imgBasicRoutingModule) ReturnBinary() {
	returnBinary(o.w, o.r, o.mime, o.size)
}

func (o *imgBasicRoutingModule) Close() {
	o.r.Close()
}

func newImgBasicRoutingModule(mime string) func(io.ReadCloser, string, http.ResponseWriter, int64) RoutingModule {
	return func(r io.ReadCloser, vpath string, w http.ResponseWriter, size int64) RoutingModule {
		if len(vpath) != 0 {
			http.Error(w, "Not Found File", http.StatusUnsupportedMediaType)
			r.Close()
			return nil
		}

		return &imgBasicRoutingModule{
			mime: mime,
			w:    w,
			r:    r,
			size: size,
		}
	}
}

func newImgBasicRoutingModule2(mime string) func(string, string, http.ResponseWriter) RoutingModule {
	return func(path, vpath string, w http.ResponseWriter) RoutingModule {
		// open binary file
		file, err := os.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return nil
		}

		// file size
		stat, err := file.Stat()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			file.Close()
			return nil
		}

		return newImgBasicRoutingModule(mime)(file, vpath, w, stat.Size())
	}
}

//
// basic image module
//

var bmpConf = install(&moduleConfig{
	name:     "bmp",
	exts:     []string{"bmp"},
	types:    typeModuleImage,
	routing:  newImgBasicRoutingModule("image/bmp"),
	routing2: newImgBasicRoutingModule2("image/bmp"),
})

var gifConf = install(&moduleConfig{
	name:     "gif",
	exts:     []string{"gif"},
	types:    typeModuleImage,
	routing:  newImgBasicRoutingModule("image/gif"),
	routing2: newImgBasicRoutingModule2("image/gif"),
})

var jpgConf = install(&moduleConfig{
	name:     "jpg",
	exts:     []string{"jpg", "jpeg"},
	types:    typeModuleImage,
	routing:  newImgBasicRoutingModule("image/jpeg"),
	routing2: newImgBasicRoutingModule2("image/jpeg"),
})

var pngConf = install(&moduleConfig{
	name:     "png",
	exts:     []string{"png"},
	types:    typeModuleImage,
	routing:  newImgBasicRoutingModule("image/png"),
	routing2: newImgBasicRoutingModule2("image/png"),
})
