package main

import (
	"log"
	"net/http"
	"strconv"

	"./archive"
	"./image"
	"./utility/array"
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
	json := "{" +
		"\"image\":" + array.ToJson(image.SupportType()) + "," +
		"\"archive\":" + array.ToJson(archive.SupportType()) +
		"}"
	buff := []byte(json)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(buff)))
	if _, err := w.Write(buff); err != nil {
		log.Printf("ERR: supportType: %v\n", err)
	}
}
