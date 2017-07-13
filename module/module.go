package module

import (
	"fmt"
	"io"
	"path/filepath"
)

type ModuleType int

const (
	MODULE_IMAGE ModuleType = iota
	MODULE_ARCH
)

const errNotSupport = "%s is not support this module type"

type File struct {
	Data io.ReadCloser
	Mime string
	Size int64
}

type Module struct {
	exts []string
	Type ModuleType
	// archive type
	FuncArchFiles func(r ReaderAt) []string
	FuncArchRead  func(r ReaderAt, vpath string) *File
	// image type
	FuncImageRead func(r Reader) *File
}

var modules []*Module

func GetSupportModule(path string) *Module {
	ext := filepath.Ext(path)

	for _, mod := range modules {
		for _, xt := range mod.exts {
			if xt == ext {
				return mod
			}
		}
	}

	return nil
}

func RegisterArchModule(exts []string, funcArchFiles func(r ReaderAt) []string, funcArchRead func(r ReaderAt, vpath string) *File) {
	h := &Module{
		exts:          exts,
		Type:          MODULE_ARCH,
		FuncArchFiles: funcArchFiles,
		FuncArchRead:  funcArchRead,
		FuncImageRead: dummyFuncImageRead,
	}
	modules = append(modules, h)
}

func RegisterImageModule(exts []string, funcImageRead func(r Reader) *File) {
	h := &Module{
		exts:          exts,
		Type:          MODULE_IMAGE,
		FuncArchFiles: dummyFuncArchFiles,
		FuncArchRead:  dummyFuncArchRead,
		FuncImageRead: funcImageRead,
	}
	modules = append(modules, h)
}

func SupportType() map[ModuleType][]string {
	types := map[ModuleType][]string{
		MODULE_IMAGE: make([]string, 0),
		MODULE_ARCH:  make([]string, 0),
	}

	for _, module := range modules {
		v, _ := types[module.Type]
		v = append(v, module.exts...)
		types[module.Type] = v
	}
	return types
}

func dummyFuncArchFiles(r ReaderAt) []string {
	fmt.Printf(errNotSupport, "Files")
	return nil
}

func dummyFuncArchRead(r ReaderAt, vpath string) *File {
	fmt.Printf(errNotSupport, "Read")
	return nil
}

func dummyFuncImageRead(r Reader) *File {
	fmt.Printf(errNotSupport, "Read")
	return nil
}
