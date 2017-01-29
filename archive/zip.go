package archive

import (
	"archive/zip"
	"log"
	"net/http"
	"strconv"

	"../image"
	"../utility/array"
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

	json := array.ToJson(paths)
	buff := []byte(json)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(buff)))
	if _, err := w.Write(buff); err != nil {
		log.Printf("ERR: zipInfoCallback: %v [%s]\n", err, path)
	}
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
