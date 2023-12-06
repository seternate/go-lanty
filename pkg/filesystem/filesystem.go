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

func SearchFileByName(file string, directory string) (filepaths []string, err error) {
	filepath.WalkDir(directory, func(path string, dirEntry fs.DirEntry, err error) error {
		file = filepath.FromSlash(file)
		if !dirEntry.IsDir() && strings.HasSuffix(path, file) {
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
	if err != nil {
		return err
	}

	return nil
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

func LoadFromYAMLFile(file string, v interface{}) (err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()

	err = LoadFromYAML(f, v)

	return
}
