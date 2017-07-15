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

type FileReader struct {
	io.ReadCloser
	size int64
}

func (o *FileReader) Size() int64 {
	return o.size
}

type FileReaderAt struct {
	ioReaderAtCloser
	size int64
}

func (o *FileReaderAt) Size() int64 {
	return o.size
}

func NewReader(file io.ReadCloser, size int64) (r Reader, err error) {
	r = &FileReader{file, size}
	return
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
