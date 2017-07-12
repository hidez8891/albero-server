package zip

import (
	"archive/zip"
	"strings"

	"github.com/hidez8891/albero-server/module"
)

func init() {
	module.RegisterArchModule([]string{".zip"}, files, read)
}

func files(r module.ReaderAt) []string {
	arch, err := zip.NewReader(r, r.Size())
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

func read(r module.ReaderAt, vpath string) *module.File {
	arch, err := zip.NewReader(r, r.Size())
	if err != nil {
		return nil
	}

	var zipfile *zip.File
	for _, f := range arch.File {
		path := strings.Replace(f.Name, "\\", "/", -1)

		if path == vpath {
			zipfile = f
			break
		}
	}

	if zipfile == nil {
		return nil
	}

	rr, err := zipfile.Open()
	if err != nil {
		return nil
	}

	file := &module.File{
		Data: rr,
		Mime: "",
		Size: int64(zipfile.FileHeader.UncompressedSize64),
	}

	return file
}
