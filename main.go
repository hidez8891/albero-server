package main

import (
	"log"
	"net/http"

	"./archive"
	"./image"
)

const defaultPort = "5358"

func main() {
	// support       -> support file/archive extension
	// image/info    -> image infomation (json)
	// image/data    -> image binary data
	// image/ext     -> support image extension
	// archive/info  -> archive infomation (json)
	// archive/data  -> archive's image binary data
	// archive/ext   -> support archive extension
	image.SetHttpRoute()
	archive.SetHttpRoute()

	log.Println("Listening on " + defaultPort)
	err := http.ListenAndServe(":"+defaultPort, nil)
	if err != nil {
		log.Fatal("Listen And Serve:", err)
	}
}
