package storageadapter

import (
	"bytes"
	"io"
	"net/url"
	"strings"
	"testing"

	"github.com/seternate/go-lanty/pkg/application/asset"
	"github.com/spf13/afero"
)

func TestNewFilesystemStorageAdapter(t *testing.T) {
	fs := afero.NewMemMapFs()
	adapter := NewFilesystemStorageAdapter(fs)

	if adapter == nil {
		t.Fatal("NewFilesystemStorageAdapter() returned nil")
	}

	if adapter.fs != fs {
		t.Error("FilesystemStorageAdapter.fs does not match provided filesystem")
	}

	// Verify it implements the interface
	var _ asset.StorageAdapter = adapter
}

func TestFilesystemStorageAdapter_Fetch(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(afero.Fs) (url.URL, []byte)
		expectError bool
		validate    func(t *testing.T, data io.ReadCloser, size uint64)
	}{
		{
			name: "successful fetch",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				testPath := "/test/file.txt"
				testData := []byte("hello world")
				_ = afero.WriteFile(fs, testPath, testData, 0644)
				return url.URL{Scheme: "file", Path: testPath}, testData
			},
			expectError: false,
			validate: func(t *testing.T, data io.ReadCloser, size uint64) {
				defer data.Close()
				if size != 11 {
					t.Errorf("expected size 11, got %d", size)
				}
				readData, err := io.ReadAll(data)
				if err != nil {
					t.Fatalf("failed to read data: %v", err)
				}
				if !bytes.Equal(readData, []byte("hello world")) {
					t.Errorf("data mismatch: got %q, want %q", readData, "hello world")
				}
			},
		},
		{
			name: "fetch empty file",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				testPath := "/test/empty.txt"
				testData := []byte("")
				_ = afero.WriteFile(fs, testPath, testData, 0644)
				return url.URL{Scheme: "file", Path: testPath}, testData
			},
			expectError: false,
			validate: func(t *testing.T, data io.ReadCloser, size uint64) {
				defer data.Close()
				if size != 0 {
					t.Errorf("expected size 0, got %d", size)
				}
				readData, err := io.ReadAll(data)
				if err != nil {
					t.Fatalf("failed to read data: %v", err)
				}
				if len(readData) != 0 {
					t.Errorf("expected empty data, got %q", readData)
				}
			},
		},
		{
			name: "fetch large file",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				testPath := "/test/large.bin"
				testData := make([]byte, 1024*1024) // 1MB
				for i := range testData {
					testData[i] = byte(i % 256)
				}
				_ = afero.WriteFile(fs, testPath, testData, 0644)
				return url.URL{Scheme: "file", Path: testPath}, testData
			},
			expectError: false,
			validate: func(t *testing.T, data io.ReadCloser, size uint64) {
				defer data.Close()
				expectedSize := uint64(1024 * 1024)
				if size != expectedSize {
					t.Errorf("expected size %d, got %d", expectedSize, size)
				}
			},
		},
		{
			name: "file not found",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				return url.URL{Scheme: "file", Path: "/nonexistent/file.txt"}, nil
			},
			expectError: true,
			validate:    nil,
		},
		{
			name: "invalid scheme",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				return url.URL{Scheme: "http", Path: "/test/file.txt"}, nil
			},
			expectError: true,
			validate:    nil,
		},
		{
			name: "empty scheme",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				return url.URL{Scheme: "", Path: "/test/file.txt"}, nil
			},
			expectError: true,
			validate:    nil,
		},
		{
			name: "file in subdirectory",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				testPath := "/deep/nested/path/file.txt"
				testData := []byte("nested file content")
				_ = fs.MkdirAll("/deep/nested/path", 0755)
				_ = afero.WriteFile(fs, testPath, testData, 0644)
				return url.URL{Scheme: "file", Path: testPath}, testData
			},
			expectError: false,
			validate: func(t *testing.T, data io.ReadCloser, size uint64) {
				defer data.Close()
				expectedSize := uint64(len("nested file content"))
				if size != expectedSize {
					t.Errorf("expected size %d, got %d", expectedSize, size)
				}
				readData, err := io.ReadAll(data)
				if err != nil {
					t.Fatalf("failed to read data: %v", err)
				}
				if !bytes.Equal(readData, []byte("nested file content")) {
					t.Errorf("data mismatch: got %q, want %q", readData, "nested file content")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			testURL, expectedData := tt.setup(fs)
			adapter := NewFilesystemStorageAdapter(fs)

			data, size, err := adapter.Fetch(testURL)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if data != nil {
					t.Error("expected nil data on error")
					if err := data.Close(); err != nil {
						t.Errorf("failed to close data on error: %v", err)
					}
				}
				// Verify error message mentions scheme validation if applicable
				if testURL.Scheme != "file" && err != nil {
					if !strings.Contains(err.Error(), "scheme") {
						t.Errorf("error message should mention scheme: %v", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if data == nil {
					t.Fatal("expected data, got nil")
				}
				if tt.validate != nil {
					tt.validate(t, data, size)
				} else {
					// Verify the data matches expected if validate wasn't called
					if expectedData != nil {
						readData, readErr := io.ReadAll(data)
						if readErr != nil {
							t.Fatalf("failed to read data: %v", readErr)
						}
						if !bytes.Equal(readData, expectedData) {
							t.Errorf("data mismatch: got %q, want %q", readData, expectedData)
						}
					}
					if err := data.Close(); err != nil {
						t.Errorf("failed to close data: %v", err)
					}
				}
			}
		})
	}
}

func TestFilesystemStorageAdapter_Save(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(afero.Fs) (url.URL, []byte)
		expectError bool
		validate    func(t *testing.T, fs afero.Fs, testURL url.URL, written uint64)
	}{
		{
			name: "successful save to new file",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				testPath := "/test/newfile.txt"
				testData := []byte("new file content")
				return url.URL{Scheme: "file", Path: testPath}, testData
			},
			expectError: false,
			validate: func(t *testing.T, fs afero.Fs, testURL url.URL, written uint64) {
				expectedWritten := uint64(len("new file content"))
				if written != expectedWritten {
					t.Errorf("expected written bytes %d, got %d", expectedWritten, written)
				}
				// Verify file exists
				exists, _ := afero.Exists(fs, testURL.Path)
				if !exists {
					t.Error("file was not created")
				}
				// Verify file content
				content, err := afero.ReadFile(fs, testURL.Path)
				if err != nil {
					t.Fatalf("failed to read saved file: %v", err)
				}
				if !bytes.Equal(content, []byte("new file content")) {
					t.Errorf("file content mismatch: got %q, want %q", content, "new file content")
				}
				// Verify temp file was removed
				matches, _ := afero.Glob(fs, testURL.Path+".tmp.*")
				if len(matches) > 0 {
					t.Error("temporary file was not removed")
				}
			},
		},
		{
			name: "save overwrites existing file",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				testPath := "/test/existing.txt"
				_ = afero.WriteFile(fs, testPath, []byte("old content"), 0644)
				testData := []byte("new content")
				return url.URL{Scheme: "file", Path: testPath}, testData
			},
			expectError: false,
			validate: func(t *testing.T, fs afero.Fs, testURL url.URL, written uint64) {
				if written != 11 {
					t.Errorf("expected written bytes 11, got %d", written)
				}
				content, err := afero.ReadFile(fs, testURL.Path)
				if err != nil {
					t.Fatalf("failed to read saved file: %v", err)
				}
				if !bytes.Equal(content, []byte("new content")) {
					t.Errorf("file content mismatch: got %q, want %q", content, "new content")
				}
			},
		},
		{
			name: "save creates directories",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				testPath := "/new/deep/path/file.txt"
				testData := []byte("content")
				return url.URL{Scheme: "file", Path: testPath}, testData
			},
			expectError: false,
			validate: func(t *testing.T, fs afero.Fs, testURL url.URL, written uint64) {
				// Verify directories were created
				dirExists, _ := afero.DirExists(fs, "/new/deep/path")
				if !dirExists {
					t.Error("directories were not created")
				}
				// Verify file exists
				exists, _ := afero.Exists(fs, testURL.Path)
				if !exists {
					t.Error("file was not created")
				}
			},
		},
		{
			name: "save empty file",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				testPath := "/test/empty.txt"
				testData := []byte("")
				return url.URL{Scheme: "file", Path: testPath}, testData
			},
			expectError: false,
			validate: func(t *testing.T, fs afero.Fs, testURL url.URL, written uint64) {
				if written != 0 {
					t.Errorf("expected written bytes 0, got %d", written)
				}
				exists, _ := afero.Exists(fs, testURL.Path)
				if !exists {
					t.Error("empty file was not created")
				}
			},
		},
		{
			name: "save large file",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				testPath := "/test/large.bin"
				testData := make([]byte, 512*1024) // 512KB
				for i := range testData {
					testData[i] = byte(i % 256)
				}
				return url.URL{Scheme: "file", Path: testPath}, testData
			},
			expectError: false,
			validate: func(t *testing.T, fs afero.Fs, testURL url.URL, written uint64) {
				expectedSize := uint64(512 * 1024)
				if written != expectedSize {
					t.Errorf("expected written bytes %d, got %d", expectedSize, written)
				}
			},
		},
		{
			name: "invalid scheme",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				return url.URL{Scheme: "http", Path: "/test/file.txt"}, []byte("content")
			},
			expectError: true,
			validate:    nil,
		},
		{
			name: "save to root file",
			setup: func(fs afero.Fs) (url.URL, []byte) {
				testPath := "/rootfile.txt"
				testData := []byte("root content")
				return url.URL{Scheme: "file", Path: testPath}, testData
			},
			expectError: false,
			validate: func(t *testing.T, fs afero.Fs, testURL url.URL, written uint64) {
				if written != 12 {
					t.Errorf("expected written bytes 12, got %d", written)
				}
				exists, _ := afero.Exists(fs, testURL.Path)
				if !exists {
					t.Error("root file was not created")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			testURL, testData := tt.setup(fs)
			adapter := NewFilesystemStorageAdapter(fs)

			dataReader := bytes.NewReader(testData)
			written, err := adapter.Save(testURL, dataReader)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if written != 0 {
					t.Errorf("expected written bytes 0 on error, got %d", written)
				}
				// Verify error message mentions scheme validation if applicable
				if testURL.Scheme != "file" && err != nil {
					if !strings.Contains(err.Error(), "scheme") {
						t.Errorf("error message should mention scheme: %v", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if tt.validate != nil {
					tt.validate(t, fs, testURL, written)
				}
			}
		})
	}
}

func TestFilesystemStorageAdapter_Save_AtomicWrite(t *testing.T) {
	fs := afero.NewMemMapFs()
	adapter := NewFilesystemStorageAdapter(fs)

	testPath := "/test/atomic.txt"
	testURL := url.URL{Scheme: "file", Path: testPath}
	testData := []byte("atomic content")

	written, err := adapter.Save(testURL, bytes.NewReader(testData))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the file exists with correct content
	exists, _ := afero.Exists(fs, testPath)
	if !exists {
		t.Fatal("file does not exist after save")
	}

	content, err := afero.ReadFile(fs, testPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	if !bytes.Equal(content, testData) {
		t.Errorf("file content mismatch: got %q, want %q", content, testData)
	}

	if written != uint64(len(testData)) {
		t.Errorf("written bytes mismatch: got %d, want %d", written, len(testData))
	}

	// Verify no temp files remain
	matches, _ := afero.Glob(fs, testPath+".tmp.*")
	if len(matches) > 0 {
		t.Errorf("temporary files remain: %v", matches)
	}
}

func TestFilesystemStorageAdapter_Delete(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(afero.Fs) url.URL
		expectError bool
		validate    func(t *testing.T, fs afero.Fs, testURL url.URL)
	}{
		{
			name: "successful delete",
			setup: func(fs afero.Fs) url.URL {
				testPath := "/test/file.txt"
				_ = afero.WriteFile(fs, testPath, []byte("content"), 0644)
				return url.URL{Scheme: "file", Path: testPath}
			},
			expectError: false,
			validate: func(t *testing.T, fs afero.Fs, testURL url.URL) {
				exists, _ := afero.Exists(fs, testURL.Path)
				if exists {
					t.Error("file still exists after delete")
				}
			},
		},
		{
			name: "delete non-existent file",
			setup: func(fs afero.Fs) url.URL {
				return url.URL{Scheme: "file", Path: "/nonexistent/file.txt"}
			},
			expectError: true,
			validate:    nil,
		},
		{
			name: "invalid scheme",
			setup: func(fs afero.Fs) url.URL {
				return url.URL{Scheme: "http", Path: "/test/file.txt"}
			},
			expectError: true,
			validate:    nil,
		},
		{
			name: "delete file in subdirectory",
			setup: func(fs afero.Fs) url.URL {
				testPath := "/deep/nested/path/file.txt"
				_ = fs.MkdirAll("/deep/nested/path", 0755)
				_ = afero.WriteFile(fs, testPath, []byte("content"), 0644)
				return url.URL{Scheme: "file", Path: testPath}
			},
			expectError: false,
			validate: func(t *testing.T, fs afero.Fs, testURL url.URL) {
				exists, _ := afero.Exists(fs, testURL.Path)
				if exists {
					t.Error("file still exists after delete")
				}
				// Verify directory still exists
				dirExists, _ := afero.DirExists(fs, "/deep/nested/path")
				if !dirExists {
					t.Error("directory should still exist after file delete")
				}
			},
		},
		{
			name: "empty scheme",
			setup: func(fs afero.Fs) url.URL {
				return url.URL{Scheme: "", Path: "/test/file.txt"}
			},
			expectError: true,
			validate:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := afero.NewMemMapFs()
			testURL := tt.setup(fs)
			adapter := NewFilesystemStorageAdapter(fs)

			err := adapter.Delete(testURL)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				// Verify error message mentions scheme validation if applicable
				if testURL.Scheme != "file" && err != nil {
					if !strings.Contains(err.Error(), "scheme") {
						t.Errorf("error message should mention scheme: %v", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if tt.validate != nil {
					tt.validate(t, fs, testURL)
				}
			}
		})
	}
}

func TestFilesystemStorageAdapter_Integration(t *testing.T) {
	fs := afero.NewMemMapFs()
	adapter := NewFilesystemStorageAdapter(fs)

	testPath := "/test/integration.txt"
	testURL := url.URL{Scheme: "file", Path: testPath}
	testData := []byte("integration test data")

	// Test Save
	t.Run("Save", func(t *testing.T) {
		written, err := adapter.Save(testURL, bytes.NewReader(testData))
		if err != nil {
			t.Fatalf("Save() error = %v", err)
		}
		if written != uint64(len(testData)) {
			t.Errorf("Save() written = %d, want %d", written, len(testData))
		}
	})

	// Test Fetch
	t.Run("Fetch", func(t *testing.T) {
		data, size, err := adapter.Fetch(testURL)
		if err != nil {
			t.Fatalf("Fetch() error = %v", err)
		}
		defer data.Close()

		if size != uint64(len(testData)) {
			t.Errorf("Fetch() size = %d, want %d", size, len(testData))
		}

		readData, err := io.ReadAll(data)
		if err != nil {
			t.Fatalf("failed to read fetched data: %v", err)
		}

		if !bytes.Equal(readData, testData) {
			t.Errorf("Fetch() data = %q, want %q", readData, testData)
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		err := adapter.Delete(testURL)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}

		exists, _ := afero.Exists(fs, testPath)
		if exists {
			t.Error("file still exists after Delete()")
		}
	})
}

func TestValidateScheme(t *testing.T) {
	tests := []struct {
		name        string
		testURL     url.URL
		scheme      string
		expectError bool
	}{
		{
			name:        "valid scheme",
			testURL:     url.URL{Scheme: "file", Path: "/test"},
			scheme:      "file",
			expectError: false,
		},
		{
			name:        "invalid scheme",
			testURL:     url.URL{Scheme: "http", Path: "/test"},
			scheme:      "file",
			expectError: true,
		},
		{
			name:        "empty scheme",
			testURL:     url.URL{Scheme: "", Path: "/test"},
			scheme:      "file",
			expectError: true,
		},
		{
			name:        "case sensitive",
			testURL:     url.URL{Scheme: "File", Path: "/test"},
			scheme:      "file",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateScheme(tt.testURL, tt.scheme)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if !strings.Contains(err.Error(), "scheme") {
					t.Errorf("error message should mention scheme: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
