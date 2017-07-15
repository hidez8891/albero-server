package image

import (
	"github.com/hidez8891/albero-server/module"
)

func init() {
	types := []struct {
		exts []string
		mime string
	}{
		{[]string{".bmp"}, "image/bmp"},
		{[]string{".gif"}, "image/gif"},
		{[]string{".jpg", ".jpeg"}, "image/jpeg"},
		{[]string{".png"}, "image/png"},
	}

	for _, t := range types {
		module.RegisterImageModule(t.exts, rawRead(t.mime))
	}
}

func rawRead(mime string) func(module.Reader) *module.File {
	return func(r module.Reader) *module.File {
		return &module.File{
			Data: r,
			Mime: mime,
			Size: r.Size(),
		}
	}
}
