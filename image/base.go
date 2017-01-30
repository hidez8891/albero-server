package image

import (
	"io"
	"net/http"
	"os"
	"strings"

	"../utility/array"
)

/*
 * image handler routing
 * /image/image_file_path
 */

const (
	root = "/image/"
)

type imgConfig struct {
	name         string
	exts         []string
	dataCallback func(w http.ResponseWriter, r io.Reader, size int64)
}

var (
	confs []*imgConfig
)

func SetHttpRoute() {
	http.HandleFunc(root, handler)
}

func SupportType() []string {
	exts := []string{}
	for _, conf := range confs {
		exts = append(exts, conf.exts...)
	}
	return exts
}

func WriteResponse(w http.ResponseWriter, path string, r io.Reader, size int64) {
	exti := strings.LastIndex(path, ".")
	if exti < 0 {
		http.Error(w, "No Support Type", http.StatusUnsupportedMediaType)
		return
	}
	ext := path[exti+1:]

	// search extension
	for _, conf := range confs {
		if array.IsInclude(ext, conf.exts) {
			conf.dataCallback(w, r, size)
			return
		}
	}

	// not found image type
	http.Error(w, "No Support Type", http.StatusUnsupportedMediaType)
	return
}

func install(h *imgConfig) *imgConfig {
	confs = append(confs, h)
	return h
}

func handler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len(root):]
	if filePath == "" {
		http.NotFound(w, r)
		return
	}

	exti := strings.LastIndex(filePath, ".")
	if exti < 0 {
		http.Error(w, "No Support Type", http.StatusUnsupportedMediaType)
		return
	}
	ext := filePath[exti+1:]

	for _, conf := range confs {
		if array.IsInclude(ext, conf.exts) {
			// open binary file
			file, err := os.Open(filePath)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer file.Close()

			// file size
			stat, err := file.Stat()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// call image callback
			conf.dataCallback(w, file, stat.Size())
			return
		}
	}

	// not found image type
	http.Error(w, "No Support Type", http.StatusUnsupportedMediaType)
	return
}
