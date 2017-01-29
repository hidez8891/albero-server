package archive

import (
	"archive/zip"
	"net/http"

	"../image"
	"../utility/array"
	"../utility/json"
)

var zipConf = install(&archConfig{
	name:         "zip",
	exts:         []string{"zip"},
	infoCallback: zipInfoCallback,
	dataCallback: zipDataCallback,
})

func zipInfoCallback(w http.ResponseWriter, path string) {
	r, err := zip.OpenReader(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Close()

	paths := make([]string, len(r.File))
	for i, f := range r.File {
		paths[i] = f.Name
	}

	json.WriteResponse(w, paths)
}

func zipDataCallback(w http.ResponseWriter, path, page string) {
	r, err := zip.OpenReader(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Close()

	paths := make([]string, len(r.File))
	for i, f := range r.File {
		paths[i] = f.Name
	}

	index := array.Search(paths, page)
	if index < 0 {
		http.Error(w, "Image file Not Found", http.StatusNotFound)
		return
	}

	size := r.File[index].FileInfo().Size()
	body, err := r.File[index].Open()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer body.Close()

	image.ResponseWrite(w, page, body, size)
}
