package json

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func WriteResponse(w http.ResponseWriter, data interface{}) {
	str, _ := json.Marshal(data)
	buff := []byte(str)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(buff)))
	if _, err := w.Write(buff); err != nil {
		log.Printf("ERR: WriteResponse: %v [%s]\n", err, str)
	}
}
