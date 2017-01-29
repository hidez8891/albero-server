package archive

import (
	"net/http"
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

func install(h *archConfig) *archConfig {
	confs = append(confs, h)
	return h
}

func SetHttpRoute() {
	http.HandleFunc(root, archHandler)
}

func archHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path[len(root):]
	if url == "" {
		http.NotFound(w, r)
		return
	}

	archPath := url
	archPage := ""

	pagei := strings.LastIndex(url, "/")
	if pagei >= 0 {
		archPath = url[:pagei]
		archPage = url[pagei+1:]
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
