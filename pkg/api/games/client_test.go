package games

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/seternate/go-lanty/pkg/api/models/game"
	"github.com/seternate/go-lanty/pkg/api/paths"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockClientInterface struct {
	restRequestErr        error
	restRequestResult     *http.Response
	restRequestCallCount  int
	lastRESTRequestMethod string
	lastRESTRequestPath   paths.Path
	lastRESTRequestParams paths.Params
	lastRESTRequestHeaders map[string]string
	lastRESTRequestBody   []byte
}

func newMockClientInterface() *mockClientInterface {
	return &mockClientInterface{}
}

func (m *mockClientInterface) BuildURL(path paths.Path, params paths.Params) (string, error) {
	return "", nil
}

func (m *mockClientInterface) NewRESTRequestWithContext(ctx context.Context, method string, url string, headers map[string]string, queryParams map[string]string, body io.Reader) (*http.Request, error) {
	if body != nil {
		m.lastRESTRequestBody, _ = io.ReadAll(body)
	}
	return nil, nil
}

func (m *mockClientInterface) Do(req *http.Request) (*http.Response, error) {
	return nil, nil
}

func (m *mockClientInterface) RESTRequest(ctx context.Context, method string, path paths.Path, pathParams paths.Params, headers map[string]string, queryParams map[string]string, body io.Reader) (*http.Response, error) {
	m.restRequestCallCount++
	m.lastRESTRequestMethod = method
	m.lastRESTRequestPath = path
	m.lastRESTRequestParams = pathParams
	m.lastRESTRequestHeaders = headers
	if body != nil {
		m.lastRESTRequestBody, _ = io.ReadAll(body)
	}
	if m.restRequestErr != nil {
		return nil, m.restRequestErr
	}
	return m.restRequestResult, nil
}

func (m *mockClientInterface) HTTPClient() *http.Client {
	return nil
}

func (m *mockClientInterface) Headers() map[string]string {
	return nil
}

type mockHashCalculator struct {
	hasher          hash.Hash
	calculateErr    error
	calculatedValue string
}

func newMockHashCalculator() *mockHashCalculator {
	return &mockHashCalculator{
		hasher: sha256.New(),
	}
}

func (m *mockHashCalculator) SHA256Hasher() hash.Hash {
	return m.hasher
}

func (m *mockHashCalculator) CalculateChecksum(hasher hash.Hash) (string, error) {
	if m.calculateErr != nil {
		return "", m.calculateErr
	}
	if m.calculatedValue != "" {
		return m.calculatedValue, nil
	}
	hashBytes := hasher.Sum(nil)
	return fmt.Sprintf("%x", hashBytes), nil
}

type mockMimeTypeDetector struct {
	mimeType string
	err      error
}

func newMockMimeTypeDetector() *mockMimeTypeDetector {
	return &mockMimeTypeDetector{
		mimeType: "image/png",
	}
}

func (m *mockMimeTypeDetector) DetectAndPreserveReader(reader io.Reader) (string, io.Reader, error) {
	if m.err != nil {
		return "", nil, m.err
	}
	data, _ := io.ReadAll(reader)
	return m.mimeType, bytes.NewReader(data), nil
}

func TestNew(t *testing.T) {
	t.Run("with nil dependencies", func(t *testing.T) {
		mockClient := newMockClientInterface()
		client := New(mockClient, nil, nil)
		
		assert.NotNil(t, client)
		assert.Equal(t, mockClient, client.apiClient)
		assert.NotNil(t, client.hashCalculator)
		assert.NotNil(t, client.mimeDetector)
	})

	t.Run("with provided dependencies", func(t *testing.T) {
		mockClient := newMockClientInterface()
		mockHashCalc := newMockHashCalculator()
		mockMimeDet := newMockMimeTypeDetector()
		
		client := New(mockClient, mockHashCalc, mockMimeDet)
		
		assert.NotNil(t, client)
		assert.Equal(t, mockClient, client.apiClient)
		assert.Equal(t, mockHashCalc, client.hashCalculator)
		assert.Equal(t, mockMimeDet, client.mimeDetector)
	})
}

func TestClient_GetGames(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedGames := []*gamemodel.Game{
			{Slug: "game1", Name: "Game 1"},
			{Slug: "game2", Name: "Game 2"},
		}
		body, _ := json.Marshal(expectedGames)
		
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		resp.Body.Write(body)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(bytes.NewReader(body))
		
		client := New(mockClient, nil, nil)
		games, err := client.GetGames(context.Background())
		
		require.NoError(t, err)
		assert.Len(t, games, 2)
		assert.Equal(t, "game1", games[0].Slug)
		assert.Equal(t, paths.GetGames.Method, mockClient.lastRESTRequestMethod)
	})

	t.Run("API request error", func(t *testing.T) {
		mockClient := newMockClientInterface()
		mockClient.restRequestErr = assert.AnError
		
		client := New(mockClient, nil, nil)
		games, err := client.GetGames(context.Background())
		
		require.Error(t, err)
		assert.Nil(t, games)
		assert.Contains(t, err.Error(), "failed API request")
	})
}

func TestClient_GetGame(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		expectedGame := gamemodel.Game{Slug: "test-game", Name: "Test Game"}
		body, _ := json.Marshal(expectedGame)
		
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		resp.Body.Write(body)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(bytes.NewReader(body))
		
		client := New(mockClient, nil, nil)
		game, err := client.GetGame(context.Background(), "test-game")
		
		require.NoError(t, err)
		assert.Equal(t, "test-game", game.Slug)
		assert.Equal(t, paths.GetGame, mockClient.lastRESTRequestPath)
		assert.Equal(t, paths.GameParams{Slug: "test-game"}, mockClient.lastRESTRequestParams)
	})
}

func TestClient_PutGame(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := &gamemodel.UpsertGameRequest{Name: "Test Game"}
		expectedGame := gamemodel.Game{Slug: "test-game", Name: "Test Game"}
		body, _ := json.Marshal(expectedGame)
		
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		resp.Body.Write(body)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(bytes.NewReader(body))
		
		client := New(mockClient, nil, nil)
		result, err := client.PutGame(context.Background(), "test-game", req)
		
		require.NoError(t, err)
		assert.Equal(t, "test-game", result.Slug)
		assert.Equal(t, "application/json", mockClient.lastRESTRequestHeaders["Content-Type"])
	})
}

func TestClient_DeleteGame(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(bytes.NewReader([]byte{}))
		
		client := New(mockClient, nil, nil)
		err := client.DeleteGame(context.Background(), "test-game")
		
		require.NoError(t, err)
		assert.Equal(t, paths.DeleteGame, mockClient.lastRESTRequestPath)
	})
}

func TestClient_GetIcon(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		iconData := []byte("icon data")
		hasher := sha256.New()
		hasher.Write(iconData)
		hashBytes := hasher.Sum(nil)
		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)
		digestHeader := fmt.Sprintf("sha-256=%s", hashBase64)
		
		resp := httptest.NewRecorder()
		resp.Header().Set("Content-Digest", digestHeader)
		resp.Header().Set("Content-Type", "image/png")
		resp.WriteHeader(http.StatusOK)
		resp.Body.Write(iconData)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(bytes.NewReader(iconData))
		
		client := New(mockClient, nil, nil)
		data, mimeType, err := client.GetIcon(context.Background(), "test-game")
		
		require.NoError(t, err)
		assert.Equal(t, iconData, data)
		assert.Equal(t, "image/png", mimeType)
	})

	t.Run("missing Content-Digest header", func(t *testing.T) {
		data := []byte("data")
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		resp.Body.Write(data)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(bytes.NewReader(data))
		
		client := New(mockClient, nil, nil)
		data, mimeType, err := client.GetIcon(context.Background(), "test-game")
		
		require.Error(t, err)
		assert.Nil(t, data)
		assert.Empty(t, mimeType)
		assert.Contains(t, err.Error(), "missing Content-Digest header")
	})

	t.Run("checksum mismatch", func(t *testing.T) {
		iconData := []byte("icon data")
		wrongHashBase64 := base64.StdEncoding.EncodeToString([]byte("wrong"))
		digestHeader := fmt.Sprintf("sha-256=%s", wrongHashBase64)
		
		resp := httptest.NewRecorder()
		resp.Header().Set("Content-Digest", digestHeader)
		resp.WriteHeader(http.StatusOK)
		resp.Body.Write(iconData)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(bytes.NewReader(iconData))
		
		client := New(mockClient, nil, nil)
		data, mimeType, err := client.GetIcon(context.Background(), "test-game")
		
		require.Error(t, err)
		assert.Nil(t, data)
		assert.Empty(t, mimeType)
		assert.Contains(t, err.Error(), "checksum mismatch")
	})

	t.Run("unsupported algorithm", func(t *testing.T) {
		iconData := []byte("icon data")
		digestHeader := "md5=ZDQxZDhjZDk4ZjAwYjIwNGU5ODAwOTk4ZWNmODQyN2U="
		
		resp := httptest.NewRecorder()
		resp.Header().Set("Content-Digest", digestHeader)
		resp.WriteHeader(http.StatusOK)
		resp.Body.Write(iconData)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(bytes.NewReader(iconData))
		
		client := New(mockClient, nil, nil)
		data, mimeType, err := client.GetIcon(context.Background(), "test-game")
		
		require.Error(t, err)
		assert.Nil(t, data)
		assert.Empty(t, mimeType)
		assert.Contains(t, err.Error(), "unsupported digest algorithm")
	})

	t.Run("mimetype detection when header missing", func(t *testing.T) {
		iconData := []byte("icon data")
		hasher := sha256.New()
		hasher.Write(iconData)
		hashBytes := hasher.Sum(nil)
		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)
		digestHeader := fmt.Sprintf("sha-256=%s", hashBase64)
		
		resp := httptest.NewRecorder()
		resp.Header().Set("Content-Digest", digestHeader)
		resp.WriteHeader(http.StatusOK)
		resp.Body.Write(iconData)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(bytes.NewReader(iconData))
		
		mockMimeDet := newMockMimeTypeDetector()
		mockMimeDet.mimeType = "image/jpeg"
		
		client := New(mockClient, nil, mockMimeDet)
		data, mimeType, err := client.GetIcon(context.Background(), "test-game")
		
		require.NoError(t, err)
		assert.Equal(t, iconData, data)
		assert.Equal(t, "image/jpeg", mimeType)
	})
}

func TestClient_PutIcon(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		iconData := []byte("icon data")
		hasher := sha256.New()
		hasher.Write(iconData)
		hashBytes := hasher.Sum(nil)
		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)
		expectedDigest := fmt.Sprintf("sha-256=%s", hashBase64)
		
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(bytes.NewReader([]byte{}))
		
		mockMimeDet := newMockMimeTypeDetector()
		mockMimeDet.mimeType = "image/png"
		
		client := New(mockClient, nil, mockMimeDet)
		err := client.PutIcon(context.Background(), "test-game", bytes.NewReader(iconData))
		
		require.NoError(t, err)
		assert.Equal(t, "image/png", mockClient.lastRESTRequestHeaders["Content-Type"])
		assert.Equal(t, expectedDigest, mockClient.lastRESTRequestHeaders["Content-Digest"])
		assert.Equal(t, iconData, mockClient.lastRESTRequestBody)
	})
}

func TestClient_DownloadBlobFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Create temporary directory
		tmpDir, err := os.MkdirTemp("", "test-download-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)
		
		blobData := []byte("blob data")
		hasher := sha256.New()
		hasher.Write(blobData)
		hashBytes := hasher.Sum(nil)
		hashBase64 := base64.StdEncoding.EncodeToString(hashBytes)
		digestHeader := fmt.Sprintf("sha-256=%s", hashBase64)
		
		resp := httptest.NewRecorder()
		resp.Header().Set("Content-Digest", digestHeader)
		resp.Header().Set("Content-Length", fmt.Sprintf("%d", len(blobData)))
		resp.WriteHeader(http.StatusOK)
		resp.Body.Write(blobData)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(bytes.NewReader(blobData))
		
		client := New(mockClient, nil, nil)
		prog, err := client.DownloadBlobFile(context.Background(), "test-game", tmpDir)
		
		require.NoError(t, err)
		assert.NotNil(t, prog)
		
		// Wait for download to complete
		timeout := time.After(5 * time.Second)
		for !prog.IsComplete() {
			select {
			case <-timeout:
				t.Fatal("timeout waiting for download")
			case <-time.After(10 * time.Millisecond):
			}
		}
		
		assert.NoError(t, prog.Error())
		filePath := filepath.Join(tmpDir, "test-game")
		data, err := os.ReadFile(filePath)
		require.NoError(t, err)
		assert.Equal(t, blobData, data)
	})
}

func TestClient_UploadBlobFromDirectory(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Create temporary directory with test files
		tmpDir, err := os.MkdirTemp("", "test-upload-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)
		
		testFile := filepath.Join(tmpDir, "test.txt")
		err = os.WriteFile(testFile, []byte("test content"), 0644)
		require.NoError(t, err)
		
		resp := httptest.NewRecorder()
		resp.WriteHeader(http.StatusOK)
		
		mockClient := newMockClientInterface()
		mockClient.restRequestResult = resp.Result()
		mockClient.restRequestResult.Body = io.NopCloser(bytes.NewReader([]byte{}))
		
		client := New(mockClient, nil, nil)
		prog, err := client.UploadBlobFromDirectory(context.Background(), "test-game", tmpDir)
		
		require.NoError(t, err)
		assert.NotNil(t, prog)
		
		// Wait for upload to complete
		timeout := time.After(5 * time.Second)
		for !prog.IsComplete() {
			select {
			case <-timeout:
				t.Fatal("timeout waiting for upload")
			case <-time.After(10 * time.Millisecond):
			}
		}
		
		assert.NoError(t, prog.Error())
		assert.Equal(t, "application/zip", mockClient.lastRESTRequestHeaders["Content-Type"])
		assert.NotEmpty(t, mockClient.lastRESTRequestHeaders["Content-Digest"])
	})
}

func TestZipDirectoryToWriter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Create temporary directory with test files
		tmpDir, err := os.MkdirTemp("", "test-zip-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)
		
		testFile1 := filepath.Join(tmpDir, "file1.txt")
		err = os.WriteFile(testFile1, []byte("content1"), 0644)
		require.NoError(t, err)
		
		testFile2 := filepath.Join(tmpDir, "file2.txt")
		err = os.WriteFile(testFile2, []byte("content2"), 0644)
		require.NoError(t, err)
		
		var buf bytes.Buffer
		err = zipDirectoryToWriter(tmpDir, &buf)
		require.NoError(t, err)
		
		// Verify zip content
		zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
		require.NoError(t, err)
		
		files := make(map[string]bool)
		for _, f := range zipReader.File {
			files[f.Name] = true
		}
		
		assert.True(t, files["file1.txt"])
		assert.True(t, files["file2.txt"])
	})

	t.Run("with subdirectory", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "test-zip-subdir-*")
		require.NoError(t, err)
		defer os.RemoveAll(tmpDir)
		
		subDir := filepath.Join(tmpDir, "subdir")
		err = os.Mkdir(subDir, 0755)
		require.NoError(t, err)
		
		testFile := filepath.Join(subDir, "file.txt")
		err = os.WriteFile(testFile, []byte("content"), 0644)
		require.NoError(t, err)
		
		var buf bytes.Buffer
		err = zipDirectoryToWriter(tmpDir, &buf)
		require.NoError(t, err)
		
		zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
		require.NoError(t, err)
		
		files := make(map[string]bool)
		for _, f := range zipReader.File {
			files[f.Name] = true
		}
		
		assert.True(t, files["subdir/"] || files["subdir"])
		assert.True(t, files["subdir/file.txt"])
	})
}
