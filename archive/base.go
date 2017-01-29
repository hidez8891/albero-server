package archive

import (
	"net/http"
	"net/url"
	"strings"

	"../utility/array"
)

/*
 * archive handler routing
 * /archive/archive_file_path[/image_file_path]
 */

const (
	root = "/archive/"
)

type archConfig struct {
	name         string
	exts         []string
	infoCallback func(w http.ResponseWriter, path string)
	dataCallback func(w http.ResponseWriter, path, page string)
}

var (
	confs []*archConfig
)

func SetHttpRoute() {
	http.HandleFunc(root, handler)
}

func install(h *archConfig) *archConfig {
	confs = append(confs, h)
	return h
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.RawPath
	if path == "" {
		path = r.URL.Path
	}
	path = path[len(root):]
	if path == "" {
		http.NotFound(w, r)
		return
	}

	archPath := path
	archPage := ""

	pagei := strings.LastIndex(path, "/")
	if pagei >= 0 {
		archPath = path[:pagei]
		archPage = path[pagei+1:]
	}

	var err error
	archPath, err = url.QueryUnescape(archPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	archPage, err = url.QueryUnescape(archPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exti := strings.LastIndex(archPath, ".")
	if exti < 0 {
		http.Error(w, "No Support Type", http.StatusUnsupportedMediaType)
		return
	}
	ext := archPath[exti+1:]

	for _, conf := range confs {
		if array.IsInclude(ext, conf.exts) {
			// call archive callback
			if archPage != "" {
				conf.dataCallback(w, archPath, archPage)
			} else {
				conf.infoCallback(w, archPath)
			}
			return
		}
	}

	// not found archive type
	http.Error(w, "No Support Type", http.StatusUnsupportedMediaType)
	return
}
