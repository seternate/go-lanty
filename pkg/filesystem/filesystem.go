package filesystem

import (
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func DetectContentTypeOfFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	contentBuffer := make([]byte, 512)
	_, err = file.Read(contentBuffer)
	if err != nil {
		return "", err
	}
	file.Seek(0, 0)

	return http.DetectContentType(contentBuffer), nil
}

func SearchFileByName(name string, searchDirectory string) (filepaths []string) {
	filepath.WalkDir(searchDirectory, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.Contains(d.Name(), name) {
			filepaths = append(filepaths, path)
		}
		return nil
	})

	return
}

func WriteToFile(source io.Reader, path string, mkdir bool) (err error) {
	directory := filepath.Dir(path)
	_, err = os.Stat(directory)

	if os.IsNotExist(err) && mkdir {
		err = os.MkdirAll(directory, 0775)
	}

	if err != nil {
		return
	}

	fileHandle, err := os.Create(path)
	if err != nil {
		return
	}
	defer fileHandle.Close()

	_, err = io.Copy(fileHandle, source)
	if err != nil {
		return err
	}

	return nil
}

func LoadFromYAML(r io.Reader, v interface{}) (err error) {
	decoder := yaml.NewDecoder(r)
	err = decoder.Decode(v)
	return
}

func LoadFromJSON(r io.Reader, v interface{}) (err error) {
	decoder := json.NewDecoder(r)
	err = decoder.Decode(v)
	return
}

func LoadFromYAMLFile(f string, v interface{}) (err error) {
	file, err := os.Open(f)
	defer file.Close()
	if err != nil {
		return
	}
	err = LoadFromYAML(file, v)
	return
}
