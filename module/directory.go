package module

import "net/http"

type directoryRoutingModule struct {
	path string
	w    http.ResponseWriter
}

func (o *directoryRoutingModule) ReturnFiles() {
	// ToDo
}

func (o *directoryRoutingModule) ReturnBinary() {
	http.Error(o.w, "Not Support", http.StatusUnsupportedMediaType)
}

func newDirectoryRouting(path string, w http.ResponseWriter) RoutingModule {
	return &directoryRoutingModule{
		path: path,
		w:    w,
	}
}
