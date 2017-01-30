package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"./archive"
	"./image"
	"./utility/json"
)

const defaultPort = 5358

func main() {
	var port uint
	flag.UintVar(&port, "port", defaultPort, "listen port")
	flag.UintVar(&port, "p", defaultPort, "listen port")
	flag.Parse()

	// image data
	// image/image_file_path
	image.SetHttpRoute()

	// archive data
	// archive/archive_file_path[/image_file_path]
	archive.SetHttpRoute()

	// support type
	http.HandleFunc("/support", supportType)

	log.Printf("Listening on %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", defaultPort), nil)
	if err != nil {
		log.Fatal("Listen And Serve:", err)
	}
}

func supportType(w http.ResponseWriter, r *http.Request) {
	data := map[string][]string{
		"image":   image.SupportType(),
		"archive": archive.SupportType(),
	}

	json.WriteResponse(w, data)
}
