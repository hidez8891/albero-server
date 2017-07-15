package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/hidez8891/albero-server/module"
	"github.com/hidez8891/albero-server/utility/array"
	"github.com/hidez8891/albero-server/utility/json"

	_ "github.com/hidez8891/albero-server/module/image"
	_ "github.com/hidez8891/albero-server/module/zip"
)

const (
	defaultPort = 5358
)

func main() {
	var port uint
	flag.UintVar(&port, "port", defaultPort, "listen port")
	flag.UintVar(&port, "p", defaultPort, "listen port")
	flag.Parse()

	// path = url.encode(path1/path2/arch/image)

	// /support return support file type (json format)
	http.HandleFunc("/support", supportRouting)

	// /type?path=path_enc return file types (json format)
	http.HandleFunc("/type", typeRouting)

	// /files?path=path_enc return files (json format)
	http.HandleFunc("/files", filesRouting)

	// /image?path=path_enc return image (binary format)
	http.HandleFunc("/image", imageRouting)

	log.Printf("Listening on %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", defaultPort), nil)
	if err != nil {
		log.Fatal("Listen And Serve:", err)
	}
}

func queryPath(r *http.Request) (string, error) {
	param := r.URL.Query()
	pathEnc := param.Get("path")
	if len(pathEnc) == 0 {
		return "", fmt.Errorf("need 'path' query")
	}

	path, err := url.QueryUnescape(pathEnc)
	if err != nil {
		return "", err
	}
	return path, nil
}

func supportRouting(w http.ResponseWriter, r *http.Request) {
	exts := struct {
		Image []string `json:"image"`
		Arch  []string `json:"arch"`
	}{}

	types := module.SupportType()
	exts.Image, _ = types[module.MODULE_IMAGE]
	exts.Arch, _ = types[module.MODULE_ARCH]

	json.WriteResponse(w, exts)
}

func typeRouting(w http.ResponseWriter, r *http.Request) {
	path, err := queryPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO path validation
	ext := filepath.Ext(path)
	supports := module.SupportType()

	ret := struct {
		Type string `json:"type"`
	}{}

	types := []struct {
		mtype module.ModuleType
		stype string
	}{
		{module.MODULE_IMAGE, "image"},
		{module.MODULE_ARCH, "arch"},
	}

	for _, tp := range types {
		exts, _ := supports[tp.mtype]
		for _, ex := range exts {
			if ex == ext {
				// return file's type string
				ret.Type = tp.stype
				json.WriteResponse(w, ret)
				return
			}
		}
	}

	// not found
	http.NotFound(w, r)
	return
}

func filesRouting(w http.ResponseWriter, r *http.Request) {
	path, err := queryPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paths := strings.SplitAfter(path, "/")
	path = ""
	for len(paths) > 0 {
		pt := filepath.Join(path, paths[0])
		_, err := os.Stat(pt)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		path = pt
		paths = paths[1:]
	}

	stat, _ := os.Stat(path)
	var files []string

	// 'paths' is not empty, 'path' needs archive file path
	// 'paths' is empty, 'path' needs archive file or directory path
	if len(paths) == 0 && stat.IsDir() {
		// directory
		infos, err := ioutil.ReadDir(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		files = make([]string, 0)
		for _, v := range infos {
			name := strings.Replace(v.Name(), "\\", "/", -1)
			if ind := strings.Index(name, "/"); ind >= 0 {
				name = name[ind+1:]
			}
			files = append(files, name)
		}
	} else {
		// archive
		mod := module.GetSupportModule(path)
		if mod == nil {
			http.Error(w, "no support type", http.StatusBadRequest)
			return
		}
		if mod.Type != module.MODULE_ARCH {
			http.Error(w, "wrong module operation", http.StatusBadRequest)
			return
		}

		file, err := module.NewReaderAt(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// TODO
		// support nested archive's path
		files = mod.FuncArchFiles(file)
		for _, vpath := range paths {
			files = array.Filter(files, func(s string) bool {
				return strings.HasPrefix(s, vpath)
			})
			files = array.Map(files, func(s string) string {
				return s[len(vpath):]
			})
		}
		files = array.Uniq(array.Map(files, func(s string) string {
			if index := strings.Index(s, "/"); index >= 0 {
				return s[:index] // remove suffix
			}
			return s
		}))
	}

	// return json format
	ret := struct {
		Files []string `json:"files"`
	}{files}
	json.WriteResponse(w, ret)
}

func imageRouting(w http.ResponseWriter, r *http.Request) {
	path, err := queryPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paths := strings.SplitAfter(path, "/")
	path = ""
	for len(paths) > 0 {
		pt := filepath.Join(path, paths[0])
		_, err := os.Stat(pt)
		if os.IsNotExist(err) {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		path = pt
		paths = paths[1:]
	}

	stat, _ := os.Stat(path)

	// 'paths' is not empty, 'path' needs archive file path
	// 'paths' is empty, 'path' needs image file
	if len(paths) == 0 {
		// image file
		// reject directory path
		if stat.IsDir() {
			http.Error(w, "path is directory", http.StatusBadRequest)
			return
		}
		imageRoutingResponse(w, path, nil)
	} else {
		// archive file
		mod := module.GetSupportModule(path)
		if mod == nil {
			http.Error(w, "no support type", http.StatusBadRequest)
			return
		}
		if mod.Type != module.MODULE_ARCH {
			http.Error(w, "wrong module operation", http.StatusBadRequest)
			return
		}

		arch, err := module.NewReaderAt(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer arch.Close()

		// TODO
		// support nested archive's path
		imgpath := strings.Join(paths, "")
		file := mod.FuncArchRead(arch, imgpath)
		if file == nil {
			http.Error(w, "fail read file", http.StatusInternalServerError)
			return
		}
		defer file.Data.Close()

		r, err := module.NewReader(file.Data, file.Size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		imageRoutingResponse(w, imgpath, r)
	}
}

func imageRoutingResponse(w http.ResponseWriter, path string, r module.Reader) {
	mod := module.GetSupportModule(path)
	if mod == nil {
		http.Error(w, "no support type", http.StatusBadRequest)
		return
	}
	if mod.Type != module.MODULE_IMAGE {
		http.Error(w, "wrong module operation", http.StatusBadRequest)
		return
	}

	if r == nil {
		var err error
		r, err = module.NewReaderAt(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer r.Close()
	}

	file := mod.FuncImageRead(r)
	if file == nil {
		http.Error(w, "fail read file", http.StatusInternalServerError)
		return
	}
	defer file.Data.Close()

	w.Header().Set("Content-Type", file.Mime)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))
	if _, err := io.Copy(w, file.Data); err != nil {
		log.Printf("ERR: WriteResponse: %v\n", err)
	}
}
