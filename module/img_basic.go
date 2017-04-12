package module

import (
	"net/http"
	"os"
)

type imgBasicRoutingModule struct {
	mime string
	path string
	w    http.ResponseWriter
}

func (o *imgBasicRoutingModule) ReturnFiles() {
	http.Error(o.w, "Not Support", http.StatusUnsupportedMediaType)
}

func (o *imgBasicRoutingModule) ReturnBinary() {
	// open binary file
	file, err := os.Open(o.path)
	if err != nil {
		http.Error(o.w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// file size
	stat, err := file.Stat()
	if err != nil {
		http.Error(o.w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return binary
	returnBinary(o.w, file, o.mime, stat.Size())
}

func newImgBasicRoutingModule(mime string) func(string, string, http.ResponseWriter) RoutingModule {
	return func(path, vpath string, w http.ResponseWriter) RoutingModule {
		if len(vpath) != 0 {
			http.Error(w, "Not Found File", http.StatusUnsupportedMediaType)
			return nil
		}

		return &imgBasicRoutingModule{
			mime: mime,
			path: path,
			w:    w,
		}
	}
}

//
// basic image module
//

var bmpConf = install(&moduleConfig{
	name:    "bmp",
	exts:    []string{"bmp"},
	routing: newImgBasicRoutingModule("image/bmp"),
})

var gifConf = install(&moduleConfig{
	name:    "gif",
	exts:    []string{"gif"},
	routing: newImgBasicRoutingModule("image/gif"),
})

var jpgConf = install(&moduleConfig{
	name:    "jpg",
	exts:    []string{"jpg", "jpeg"},
	routing: newImgBasicRoutingModule("image/jpeg"),
})

var pngConf = install(&moduleConfig{
	name:    "png",
	exts:    []string{"png"},
	routing: newImgBasicRoutingModule("image/png"),
})
