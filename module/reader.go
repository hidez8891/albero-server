package module

type Reader interface {
	Close() error
	Read(p []byte) (n int, err error)
	Size() int64
}

type ReaderAt interface {
	Reader
	ReadAt(p []byte, off int64) (n int, err error)
}
