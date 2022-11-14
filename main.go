package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var version = flag.String("v", "v1", "v1")

func main() {
	router := gin.Default()

	router.GET("", func(c *gin.Context) {
		flag.Parse()
		hostname, _ := os.Hostname()

		c.String(http.StatusOK, "This is version:%s running in pod %s", *version, hostname)
	})

	router.Run(":0317")
}

type Config struct {
	Path string
}

type FilePathTool struct {
	conf Config
}

type File struct {
	Name string
}

type Folder struct {
	Name    string
	Files   []File
	Folders map[string]*Folder
}

func newFolder(name string) *Folder {
	return &Folder{name, []File{}, make(map[string]*Folder)}
}

func (f *Folder) getFolder(name string) *Folder {
	if nextF, ok := f.Folders[name]; ok {
		return nextF
	} else {
		log.Println("Expected nested folder in", name, f.Name)
		newFolder(name)
	}
	return &Folder{} // cannot happen
}

func (f *Folder) addFolder(path []string) {
	for i, segment := range path {
		// if i == 0 {
		// 	continue
		// }
		if i == len(path)-1 { // last segment == new folder
			f.Folders[segment] = newFolder(segment)
		} else {
			f.getFolder(segment).addFolder(path[1:])
		}
	}
}

func (f *Folder) addFile(path []string) {
	for i, segment := range path {
		// if i == 0 {
		// 	continue
		// }
		if i == len(path)-1 { // last segment == file
			f.Files = append(f.Files, File{segment})
		} else {
			f.getFolder(segment).addFile(path[1:])
			return
		}
	}
}

func (f *Folder) String() string {
	var str string
	for _, file := range f.Files {
		str += f.Name + string(filepath.Separator) + file.Name + "\n"
	}
	for _, folder := range f.Folders {
		str += folder.String()
	}
	return str
}

type FilePathDetail struct {
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Type      bool              `json:"type"`
	Childrens []*FilePathDetail `json:"childrens"`
}

func (f *Folder) Get() error {
	startPath, _ := filepath.Abs(`D:\yp\my_work_dir\File-path\`)
	index := 0
	fd := &FilePathDetail{
		Childrens: make([]*FilePathDetail, 0),
	}
	visit := func(path string, info os.FileInfo, err error) error {
		index++
		fmt.Println("index", index)
		if info.IsDir() {
			fmt.Println("dir_name", info.Name(), path)
			fd.add(info, path)
		} else {
			fd.add(info, path)
			fmt.Println("file_name", info.Name(), path)
		}
		return nil
	}

	err := filepath.Walk(startPath, visit)
	if err != nil {
		return err
	}

	log.Printf("%+v\n", fd)
	return nil
}

func (fd *FilePathDetail) add(fileInfo os.FileInfo, path string) *FilePathDetail {
	if fd.Path == "" {
		fd.Name = fileInfo.Name()
		fd.Path = path
		fd.Type = fileInfo.IsDir()
		fd.Childrens = make([]*FilePathDetail, 0)
		return fd
	}

	if strings.Contains(path, fd.Path) {
		if fileInfo.IsDir() {
			f := &FilePathDetail{}
			fmt.Println("存在子集，继续往下添加")
			f.add(fileInfo, path)
			fd.Childrens = append(fd.Childrens, f)
			return f
		}
		fd.Name = fileInfo.Name()
		fd.Path = path
		fd.Type = fileInfo.IsDir()
		fd.Childrens = make([]*FilePathDetail, 0)
		return fd

	}
	// fd.Name = fileInfo.Name()
	// fd.Path = path
	// fd.Type = fileInfo.IsDir()
	// fd.Childrens = make([]*FilePathDetail, 0)
	return fd
}
