package filesystem

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func SearchFilesBreadthFirst(base string, file string, depth int, n int) (paths []string, err error) {
	if n == 0 {
		return
	}

	path := filepath.Join(base, file)
	if _, internalErr := os.Stat(path); !errors.Is(internalErr, os.ErrNotExist) {
		pathAbsolute, err := filepath.Abs(path)
		paths = append(paths, pathAbsolute)
		if err != nil {
			return paths, err
		}
	}

	if depth == 0 || (len(paths) >= n && n != -1) {
		return
	}

	nodes := make([][]string, depth+1)
	nodes[0] = []string{base}
	for nodedepth := 1; nodedepth <= depth; nodedepth++ {

		fmt.Println(nodedepth, depth)

		for _, node := range nodes[nodedepth-1] {
			childnodes, err := os.ReadDir(node)
			if err != nil {
				return paths, err
			}
			for _, childnode := range childnodes {
				if childnode.IsDir() {
					childnodeInfo, err := childnode.Info()
					if err != nil {
						return paths, err
					}
					childnodepath := filepath.Join(node, childnodeInfo.Name())
					nodes[nodedepth] = append(nodes[nodedepth], childnodepath)
				}
			}
		}

		if len(nodes[nodedepth]) == 0 {
			fmt.Printf("%s - %d", "nomorenodes", nodedepth)
			break
		}

		for _, node := range nodes[nodedepth] {
			if len(paths) >= n && n != -1 {
				return
			}
			path := filepath.Join(node, file)
			if _, internalErr := os.Stat(path); !errors.Is(internalErr, os.ErrNotExist) {
				pathAbsolute, err := filepath.Abs(path)
				paths = append(paths, pathAbsolute)
				if err != nil {
					return paths, err
				}
			}
		}
	}

	return
}

func SearchAllFilesDepthFirst() {

}

func DetectContentTypeOfFile(path string) (contenttype string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	contentBuffer := make([]byte, 512)
	_, err = file.Read(contentBuffer)
	if err != nil {
		return
	}
	contenttype = http.DetectContentType(contentBuffer)
	return
}

func SearchFileByNameLazy(file string, directory string) (filepaths []string, err error) {
	filepath.WalkDir(directory, func(path string, dirEntry fs.DirEntry, err error) error {
		file = filepath.FromSlash(file)
		if !dirEntry.IsDir() && (strings.Split(dirEntry.Name(), ".")[0] == file || strings.Compare(dirEntry.Name(), file) == 0) {
			filepaths = append(filepaths, path)
		}
		return nil
	})

	if len(filepaths) == 0 {
		err = errors.New("no file found")
	}

	return
}

func WriteToFile(source io.Reader, path string, mkdir bool) (err error) {
	directory := filepath.Dir(path)
	if _, err = os.Stat(directory); os.IsNotExist(err) && mkdir {
		err = os.MkdirAll(directory, 0775)
		if err != nil {
			return
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, source)
	return
}

func LoadFromYAML(reader io.Reader, v interface{}) (err error) {
	decoder := yaml.NewDecoder(reader)
	err = decoder.Decode(v)
	return
}

func LoadFromJSON(reader io.Reader, v interface{}) (err error) {
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(v)
	return
}

func LoadFromYAMLFile(path string, v interface{}) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	err = LoadFromYAML(file, v)
	return
}

func SaveToYAMLFile(path string, v interface{}) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	err = encoder.Encode(v)
	return
}
