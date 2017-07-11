package module

import (
	"fmt"
	"io"
	"strings"
)

type ModuleType int

const (
	MODULE_IMAGE ModuleType = iota
	MODULE_ARCH
)

const errNotSupport = "%s is not support this module type"

type File struct {
	Data []byte
	Mime string
	Size int64
}

type Module struct {
	exts []string
	Type ModuleType
	// archive type
	funcArchFiles func(r io.Reader) []string
	funcArchRead  func(r io.Reader, vpath string) *File
	// image type
	funcImageRead func(r io.Reader) *File
}

var modules []*Module

func GetSupportModule(path string) *Module {
	ext := getFileExt(path)

	for _, mod := range modules {
		for _, xt := range mod.exts {
			if xt == ext {
				return mod
			}
		}
	}

	return nil
}

func RegisterArchModule(exts []string, funcArchFiles func(r io.Reader) []string, funcArchRead func(r io.Reader, vpath string) *File) {
	h := &Module{
		exts:          exts,
		Type:          MODULE_ARCH,
		funcArchFiles: funcArchFiles,
		funcArchRead:  funcArchRead,
		funcImageRead: dummyFuncImageRead,
	}
	modules = append(modules, h)
}

func RegisterImageModule(exts []string, funcImageRead func(r io.Reader) *File) {
	h := &Module{
		exts:          exts,
		Type:          MODULE_IMAGE,
		funcArchFiles: dummyFuncArchFiles,
		funcArchRead:  dummyFuncArchRead,
		funcImageRead: funcImageRead,
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

func dummyFuncArchFiles(r io.Reader) []string {
	fmt.Printf(errNotSupport, "Files")
	return nil
}

func dummyFuncArchRead(r io.Reader, vpath string) *File {
	fmt.Printf(errNotSupport, "Read")
	return nil
}

func dummyFuncImageRead(r io.Reader) *File {
	fmt.Printf(errNotSupport, "Read")
	return nil
}

func getFileExt(path string) (ext string) {
	index := strings.LastIndex(path, ".")
	if index < 0 || index >= len(path)-1 {
		return
	}

	ext = path[index+1:]
	return
}
