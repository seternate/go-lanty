package storageadapter

import (
	"fmt"
	"io"
	"net/url"
	"path/filepath"

	"github.com/seternate/go-lanty/pkg/application/asset"
	"github.com/spf13/afero"
)

var _ asset.StorageAdapter = (*FilesystemStorageAdapter)(nil)

type FilesystemStorageAdapter struct {
	fs afero.Fs
}

func NewFilesystemStorageAdapter(fs afero.Fs) *FilesystemStorageAdapter {
	return &FilesystemStorageAdapter{
		fs: fs,
	}
}

func (adapter *FilesystemStorageAdapter) Fetch(url url.URL) (data io.ReadCloser, size uint64, err error) {
	err = validateScheme(url, "file")
	if err != nil {
		return nil, 0, fmt.Errorf("failed validating scheme: %w", err)
	}

	file, err := adapter.fs.Open(url.Path)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to open file for path=%s: %w", url.Path, err)
	}

	fileinfo, err := file.Stat()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to stat file for path=%s: %w", url.Path, err)
	}

	return file, uint64(fileinfo.Size()), nil
}

func (adapter *FilesystemStorageAdapter) Save(url url.URL, data io.Reader) (written uint64, err error) {
	err = validateScheme(url, "file")
	if err != nil {
		return 0, fmt.Errorf("failed validating scheme for url=%s: %w", url.String(), err)
	}

	dir := filepath.Dir(url.Path)
	err = adapter.fs.MkdirAll(dir, 0755)
	if err != nil {
		return 0, fmt.Errorf("failed to create directory for dir=%s: %w", dir, err)
	}

	tempFile, err := afero.TempFile(adapter.fs, dir, filepath.Base(url.Path)+".tmp.")
	if err != nil {
		return 0, fmt.Errorf("failed to create temporary file for dir=%s: %w", dir, err)
	}

	tempPath := tempFile.Name()
	defer func() {
		tempFile.Close()
		if err != nil {
			adapter.fs.Remove(tempPath)
		}
	}()

	bytesWritten, err := io.Copy(tempFile, data)
	if err != nil {
		return 0, fmt.Errorf("failed to write data to temporary file for path=%s: %w", tempPath, err)
	}

	err = tempFile.Close()
	if err != nil {
		return 0, fmt.Errorf("failed to close temporary file for path=%s: %w", tempPath, err)
	}

	err = adapter.fs.Rename(tempPath, url.Path)
	if err != nil {
		return 0, fmt.Errorf("failed to rename temporary file for temp_path=%s dest_path=%s: %w", tempPath, url.Path, err)
	}

	return uint64(bytesWritten), nil
}

func (adapter *FilesystemStorageAdapter) Delete(url url.URL) error {
	err := validateScheme(url, "file")
	if err != nil {
		return fmt.Errorf("failed validating scheme: %w", err)
	}

	err = adapter.fs.Remove(url.Path)
	if err != nil {
		return fmt.Errorf("failed to delete file for path=%s: %w", url.Path, err)
	}
	return nil
}

func validateScheme(url url.URL, scheme string) error {
	if url.Scheme != scheme {
		return fmt.Errorf("unsupported scheme: got %s - want %s", url.Scheme, scheme)
	}
	return nil
}
