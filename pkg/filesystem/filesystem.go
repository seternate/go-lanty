package filesystem

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

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
		if !dirEntry.IsDir() && strings.Split(dirEntry.Name(), ".")[0] == file {
			filepaths = append(filepaths, path)
		}
		return nil
	})

	if len(filepaths) == 0 {
		err = errors.New("no file found")
	}

	return
}

func SearchFileByName(directory string, file string, depth int) (absolutePath string, err error) {
	absolutePath, err = GetAbsolutePath(filepath.Join(directory, file))
	if err == nil {
		return
	}
	if depth == 0 {
		err = os.ErrNotExist
		return
	}
	directories := make([][]string, depth)
	directories[0] = []string{directory}
	for index := 0; index < depth; index++ {
		for _, dir := range directories[index] {
			absolutePath, childdirectories, err := searchFileByNameInChildDirectories(dir, file)
			if index < depth-1 {
				directories[index+1] = append(directories[index+1], childdirectories...)
			}
			if err == nil {
				return absolutePath, nil
			}
		}
	}
	err = os.ErrNotExist
	return
}

func GetAbsolutePath(file string) (absolutePath string, err error) {
	_, err = os.Stat(file)
	if !errors.Is(err, os.ErrNotExist) {
		absolutePath, err = filepath.Abs(file)
	}
	return
}

func GetDirectories(directory string) (directories []string, err error) {
	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		return
	}
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			continue
		}
		info, err := dirEntry.Info()
		if err != nil {
			return directories, err
		}
		directories = append(directories, filepath.Join(directory, info.Name()))
	}
	return directories, nil
}

func searchFileByNameInChildDirectories(directory string, file string) (absolutePath string, childdirectories []string, err error) {
	childdirectories, err = GetDirectories(directory)
	if err != nil {
		return
	}
	for _, childdirectory := range childdirectories {
		path := filepath.Join(childdirectory, file)
		absolutePath, err = GetAbsolutePath(path)
		if err == nil {
			return
		}
	}
	err = os.ErrNotExist
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
