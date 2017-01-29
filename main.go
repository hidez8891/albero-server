package main

import (
	"log"
	"net/http"

	"./archive"
	"./image"
	"./utility/json"
)

const defaultPort = "5358"

func main() {
	// image data
	// image/image_file_path
	image.SetHttpRoute()

	// archive data
	// archive/archive_file_path[/image_file_path]
	archive.SetHttpRoute()

	// support type
	http.HandleFunc("/support", supportType)

	log.Println("Listening on " + defaultPort)
	err := http.ListenAndServe(":"+defaultPort, nil)
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
