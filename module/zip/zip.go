package zip

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"strings"

	"github.com/hidez8891/albero-server/module"
)

func init() {
	module.RegisterArchModule([]string{"zip"}, files, read)
}

func files(r io.Reader) []string {
	buff, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}

	arch, err := zip.NewReader(bytes.NewReader(buff), int64(len(buff)))
	if err != nil {
		return nil
	}

	files := make([]string, 0)
	for _, f := range arch.File {
		path := strings.Replace(f.Name, "\\", "/", -1)
		files = append(files, path)
	}

	return files
}

func read(r io.Reader, vpath string) *module.File {
	buff, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}

	arch, err := zip.NewReader(bytes.NewReader(buff), int64(len(buff)))
	if err != nil {
		return nil
	}

	data := make([]byte, 0)
	for _, f := range arch.File {
		path := strings.Replace(f.Name, "\\", "/", -1)

		if path == vpath {
			rr, err := f.Open()
			if err != nil {
				return nil
			}
			defer rr.Close()

			data, err = ioutil.ReadAll(rr)
			if err != nil {
				return nil
			}

			break
		}
	}

	if len(data) == 0 {
		return nil
	}

	file := &module.File{
		Data: data,
		Mime: "",
		Size: int64(len(data)),
	}

	return file
}
