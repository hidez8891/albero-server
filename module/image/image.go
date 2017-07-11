package image

import (
	"io"
	"io/ioutil"

	"github.com/hidez8891/albero-server/module"
)

func init() {
	types := []struct {
		exts []string
		mime string
	}{
		{[]string{"bmp"}, "image/bmp"},
		{[]string{"gif"}, "image/gif"},
		{[]string{"jpg", "jpeg"}, "image/jpeg"},
		{[]string{"png"}, "image/png"},
	}

	for _, t := range types {
		module.RegisterImageModule(t.exts, rawRead(t.mime))
	}
}

func rawRead(mime string) func(io.Reader) *module.File {
	return func(r io.Reader) *module.File {
		buff, err := ioutil.ReadAll(r)
		if err != nil {
			return nil
		}

		return &module.File{
			Data: buff,
			Mime: mime,
			Size: int64(len(buff)),
		}
	}
}
