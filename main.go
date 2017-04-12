package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"./module"
	"./utility/json"
)

const (
	defaultPort = 5358
	rootFs      = "/fs/"
	rootImg     = "/img/"
)

func main() {
	var port uint
	flag.UintVar(&port, "port", defaultPort, "listen port")
	flag.UintVar(&port, "p", defaultPort, "listen port")
	flag.Parse()

	// fs :: get path files :: json
	http.HandleFunc(rootFs, fsRouting)

	// img :: get image binary :: binary
	http.HandleFunc(rootImg, imgRouting)

	// read support type :: json
	http.HandleFunc("/support", supportType)

	log.Printf("Listening on %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", defaultPort), nil)
	if err != nil {
		log.Fatal("Listen And Serve:", err)
	}
}

func supportType(w http.ResponseWriter, r *http.Request) {
	data := map[string][]string{
	//"image":   image.SupportType(),
	//"archive": archive.SupportType(),
	}

	json.WriteResponse(w, data)
}

func fsRouting(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(rootFs):]
	if path == "" {
		http.NotFound(w, r)
		return
	}

	_, err := os.Stat(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	router := module.Routing(path, w)
	if router != nil {
		router.ReturnFiles()
	}
}

func imgRouting(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(rootImg):]
	if path == "" {
		http.NotFound(w, r)
		return
	}

	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	router := module.Routing(path, w)
	if router != nil {
		router.ReturnBinary()
	}
}
