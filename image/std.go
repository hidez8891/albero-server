package image

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

var bmpConf = install(&imgConfig{
	name:         "bmp",
	exts:         []string{"bmp"},
	dataCallback: genDataCallback("image/bmp"),
})

var gifConf = install(&imgConfig{
	name:         "gif",
	exts:         []string{"gif"},
	dataCallback: genDataCallback("image/gif"),
})

var jpgConf = install(&imgConfig{
	name:         "jpg",
	exts:         []string{"jpg", "jpeg"},
	dataCallback: genDataCallback("image/jpeg"),
})

var pngConf = install(&imgConfig{
	name:         "png",
	exts:         []string{"png"},
	dataCallback: genDataCallback("image/png"),
})

func genDataCallback(mime string) func(http.ResponseWriter, io.Reader, int64) {
	return func(w http.ResponseWriter, r io.Reader, size int64) {
		w.Header().Set("Content-Type", mime)
		w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
		if _, err := io.Copy(w, r); err != nil {
			log.Printf("ERR: genDataReaderCallback[%s]: %v\n", mime, err)
		}
	}
}
