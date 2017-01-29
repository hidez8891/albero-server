package main

import (
	"log"
	"net/http"

	"./archive"
	"./image"
)

const defaultPort = "5358"

func main() {
	// image/image_file_path
	image.SetHttpRoute()

	// archive/archive_file_path[/image_file_path]
	archive.SetHttpRoute()

	log.Println("Listening on " + defaultPort)
	err := http.ListenAndServe(":"+defaultPort, nil)
	if err != nil {
		log.Fatal("Listen And Serve:", err)
	}
}
