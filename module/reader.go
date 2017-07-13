package module

import (
	"io"
	"os"
)

type ioReaderAtCloser interface {
	io.ReadCloser
	io.ReaderAt
}

type Reader interface {
	Close() error
	Read(p []byte) (n int, err error)
	Size() int64
}

type ReaderAt interface {
	Reader
	ReadAt(p []byte, off int64) (n int, err error)
}

type FileReaderAt struct {
	ioReaderAtCloser
	size int64
}

func (o *FileReaderAt) Size() int64 {
	return o.size
}

func NewReaderAt(path string) (r ReaderAt, err error) {
	stat, err := os.Stat(path)
	if err != nil {
		return
	}

	file, err := os.Open(path)
	if err != nil {
		return
	}

	r = &FileReaderAt{file, stat.Size()}
	return
}
