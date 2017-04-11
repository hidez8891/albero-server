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

	// read directory :: json
	// localhost/fs/base64(root/dir_path)

	// read archive :: json
	// localhost/fs/base64(root/arch.ext)

	// read archive inner directory :: json
	// localhost/fs/base64(root/arch.ext/dir_path)

	// read image :: binary
	// localhost/img/base64(root/image_file_path)

	// read image into archive :: binary
	// localhost/img/base64(root/arch.ext/image_file_path)

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
		"image":   image.SupportType(),
		"archive": archive.SupportType(),
	}

	json.WriteResponse(w, data)
}
