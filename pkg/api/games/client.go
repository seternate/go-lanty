package games

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/seternate/go-lanty/internal/interface/http/adapter/header"
	"github.com/seternate/go-lanty/pkg/api/content"
	"github.com/seternate/go-lanty/pkg/api/internal"
	"github.com/seternate/go-lanty/pkg/api/models/game"
	"github.com/seternate/go-lanty/pkg/api/paths"
	"github.com/seternate/go-lanty/pkg/api/progress"
)

type Client struct {
	apiClient      internal.ClientInterface
	hashCalculator content.HashCalculator
	mimeDetector   content.MimeTypeDetector
}

func New(apiClient internal.ClientInterface, hashCalculator content.HashCalculator, mimeDetector content.MimeTypeDetector) *Client {
	if hashCalculator == nil {
		hashCalculator = content.NewDefaultHashCalculator()
	}
	if mimeDetector == nil {
		mimeDetector = content.NewDefaultMimeTypeDetector()
	}
	return &Client{
		apiClient:      apiClient,
		hashCalculator: hashCalculator,
		mimeDetector:   mimeDetector,
	}
}

func (c *Client) GetGames(ctx context.Context) ([]*gamemodel.Game, error) {
	resp, err := c.apiClient.RESTRequest(ctx, paths.GetGames.Method, paths.GetGames, nil, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed API request: %w", err)
	}
	defer resp.Body.Close()

	var games []*gamemodel.Game
	if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return games, nil
}

func (c *Client) GetGame(ctx context.Context, slug string) (*gamemodel.Game, error) {
	resp, err := c.apiClient.RESTRequest(ctx, paths.GetGame.Method, paths.GetGame, paths.GameParams{Slug: slug}, nil, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed API request: %w", err)
	}
	defer resp.Body.Close()

	var gameResp gamemodel.Game
	if err := json.NewDecoder(resp.Body).Decode(&gameResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &gameResp, nil
}

func (c *Client) PutGame(ctx context.Context, slug string, req *gamemodel.UpsertGameRequest) (*gamemodel.Game, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.apiClient.RESTRequest(ctx, paths.PutGame.Method, paths.PutGame, paths.GameParams{Slug: slug}, map[string]string{"Content-Type": "application/json"}, nil, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed API request: %w", err)
	}
	defer resp.Body.Close()

	var gameResp gamemodel.Game
	if err := json.NewDecoder(resp.Body).Decode(&gameResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &gameResp, nil
}

func (c *Client) DeleteGame(ctx context.Context, slug string) error {
	resp, err := c.apiClient.RESTRequest(ctx, paths.DeleteGame.Method, paths.DeleteGame, paths.GameParams{Slug: slug}, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed API request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) GetIcon(ctx context.Context, slug string) ([]byte, string, error) {
	resp, err := c.apiClient.RESTRequest(ctx, paths.GetGameIcon.Method, paths.GetGameIcon, paths.GameParams{Slug: slug}, nil, nil, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed API request: %w", err)
	}
	defer resp.Body.Close()

	var bodyBuf bytes.Buffer
	_, err = io.Copy(&bodyBuf, resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response body: %w", err)
	}

	digestHeader := resp.Header.Get("Content-Digest")
	if digestHeader == "" {
		return nil, "", fmt.Errorf("missing Content-Digest header in response")
	}

	algorithm, expectedChecksumHex, err := header.DecodeContentDigestHeader(digestHeader)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse Content-Digest header: %w", err)
	}

	if algorithm != "sha-256" {
		return nil, "", fmt.Errorf("unsupported digest algorithm: %s", algorithm)
	}

	hasher := c.hashCalculator.SHA256Hasher()
	_, err = io.Copy(hasher, bytes.NewReader(bodyBuf.Bytes()))
	if err != nil {
		return nil, "", fmt.Errorf("failed to calculate checksum: %w", err)
	}

	actualChecksumHex, err := c.hashCalculator.CalculateChecksum(hasher)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get checksum: %w", err)
	}

	if actualChecksumHex != expectedChecksumHex {
		return nil, "", fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksumHex, actualChecksumHex)
	}

	mimetype := resp.Header.Get("Content-Type")
	if mimetype == "" {
		detectedMimetype, _, err := c.mimeDetector.DetectAndPreserveReader(bytes.NewReader(bodyBuf.Bytes()))
		if err != nil {
			return nil, "", fmt.Errorf("failed to detect mimetype: %w", err)
		}
		mimetype = detectedMimetype
	}

	return bodyBuf.Bytes(), mimetype, nil
}

func (c *Client) PutIcon(ctx context.Context, slug string, reader io.Reader) error {
	var dataBuf bytes.Buffer
	_, err := io.Copy(&dataBuf, reader)
	if err != nil {
		return fmt.Errorf("failed to read data: %w", err)
	}

	contentType, _, err := c.mimeDetector.DetectAndPreserveReader(bytes.NewReader(dataBuf.Bytes()))
	if err != nil {
		return fmt.Errorf("failed to detect content type: %w", err)
	}

	hasher := c.hashCalculator.SHA256Hasher()
	_, err = io.Copy(hasher, bytes.NewReader(dataBuf.Bytes()))
	if err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	checksumHex, err := c.hashCalculator.CalculateChecksum(hasher)
	if err != nil {
		return fmt.Errorf("failed to get checksum: %w", err)
	}

	checksumBytes, err := hex.DecodeString(checksumHex)
	if err != nil {
		return fmt.Errorf("failed to decode hex checksum: %w", err)
	}
	checksumBase64 := base64.StdEncoding.EncodeToString(checksumBytes)
	digest := fmt.Sprintf("sha-256=%s", checksumBase64)

	resp, err := c.apiClient.RESTRequest(
		ctx,
		paths.PutGameIcon.Method,
		paths.PutGameIcon,
		paths.GameParams{Slug: slug},
		map[string]string{"Content-Type": contentType, "Content-Digest": digest},
		nil,
		bytes.NewReader(dataBuf.Bytes()),
	)
	if err != nil {
		return fmt.Errorf("failed API request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) DownloadBlobFile(ctx context.Context, slug string, directory string) (*progress.Progress, error) {
	prog := progress.NewProgress()

	resp, err := c.apiClient.RESTRequest(ctx, paths.GetGameBlob.Method, paths.GetGameBlob, paths.GameParams{Slug: slug}, nil, nil, nil)
	if err != nil {
		err = fmt.Errorf("failed API request: %w", err)
		prog.SetError(err)
		return prog, err
	}

	contentLengthStr := resp.Header.Get("Content-Length")
	if contentLengthStr != "" {
		totalBytes, err := strconv.ParseInt(contentLengthStr, 10, 64)
		if err == nil && totalBytes >= 0 {
			prog.SetTotalBytes(totalBytes)
		}
	}

	digestHeader := resp.Header.Get("Content-Digest")
	if digestHeader == "" {
		resp.Body.Close()
		err := fmt.Errorf("missing Content-Digest header in response")
		prog.SetError(err)
		return prog, err
	}

	algorithm, expectedChecksumHex, err := header.DecodeContentDigestHeader(digestHeader)
	if err != nil {
		resp.Body.Close()
		err = fmt.Errorf("failed to parse Content-Digest header: %w", err)
		prog.SetError(err)
		return prog, err
	}

	if algorithm != "sha-256" {
		resp.Body.Close()
		err := fmt.Errorf("unsupported digest algorithm: %s", algorithm)
		prog.SetError(err)
		return prog, err
	}

	filePath := filepath.Join(directory, slug)

	if err := os.MkdirAll(directory, 0755); err != nil {
		resp.Body.Close()
		err = fmt.Errorf("failed to create directory: %w", err)
		prog.SetError(err)
		return prog, err
	}

	tempFile, err := os.CreateTemp(directory, slug+".tmp.")
	if err != nil {
		resp.Body.Close()
		err = fmt.Errorf("failed to create temporary file: %w", err)
		prog.SetError(err)
		return prog, err
	}

	tempPath := tempFile.Name()

	go func() {
		defer resp.Body.Close()
		defer func() {
			if r := recover(); r != nil {
				prog.SetError(fmt.Errorf("panic during download: %v", r))
			}
		}()
		cleanup := true
		defer func() {
			tempFile.Close()
			if cleanup {
				os.Remove(tempPath)
			}
		}()

		hasher := c.hashCalculator.SHA256Hasher()

		progressReader := progress.NewReader(resp.Body, prog)

		teeReader := io.TeeReader(progressReader, hasher)

		_, err = io.Copy(tempFile, teeReader)
		if err != nil {
			prog.SetError(fmt.Errorf("failed to write to temporary file: %w", err))
			return
		}

		if err := tempFile.Close(); err != nil {
			prog.SetError(fmt.Errorf("failed to close temporary file: %w", err))
			return
		}

		actualChecksumHex, err := c.hashCalculator.CalculateChecksum(hasher)
		if err != nil {
			prog.SetError(fmt.Errorf("failed to get checksum: %w", err))
			return
		}

		if actualChecksumHex != expectedChecksumHex {
			prog.SetError(fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksumHex, actualChecksumHex))
			return
		}

		if err := os.Rename(tempPath, filePath); err != nil {
			prog.SetError(fmt.Errorf("failed to rename temporary file to destination: %w", err))
			return
		}

		prog.MarkComplete()

		cleanup = false
	}()

	return prog, nil
}

func (c *Client) UploadBlobFromDirectory(ctx context.Context, slug string, directory string) (*progress.Progress, error) {
	prog := progress.NewProgress()

	dirInfo, err := os.Stat(directory)
	if err != nil {
		err := fmt.Errorf("failed to stat directory: %w", err)
		prog.SetError(err)
		return prog, err
	}
	if !dirInfo.IsDir() {
		err := fmt.Errorf("path is not a directory: %s", directory)
		prog.SetError(err)
		return prog, err
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				prog.SetError(fmt.Errorf("panic during upload: %v", r))
			}
		}()

		hasher := c.hashCalculator.SHA256Hasher()
		err = zipDirectoryToWriter(directory, hasher)
		if err != nil {
			prog.SetError(fmt.Errorf("failed to calculate digest: %w", err))
			return
		}

		checksumHex, err := c.hashCalculator.CalculateChecksum(hasher)
		if err != nil {
			prog.SetError(fmt.Errorf("failed to get checksum: %w", err))
			return
		}

		checksumBytes, err := hex.DecodeString(checksumHex)
		if err != nil {
			prog.SetError(fmt.Errorf("failed to decode hex checksum: %w", err))
			return
		}
		checksumBase64 := base64.StdEncoding.EncodeToString(checksumBytes)
		digest := fmt.Sprintf("sha-256=%s", checksumBase64)

		contentType := "application/zip"

		pipeReader, pipeWriter := io.Pipe()

		zipErrChan := make(chan error, 1)
		go func() {
			defer close(zipErrChan)
			defer pipeWriter.Close()

			err := zipDirectoryToWriter(directory, pipeWriter)
			if err != nil {
				zipErrChan <- err
				return
			}
		}()

		progressReader := progress.NewReader(pipeReader, prog)

		uploadErrChan := make(chan error, 1)
		go func() {
			defer close(uploadErrChan)

			resp, err := c.apiClient.RESTRequest(
				ctx,
				paths.PutGameBlob.Method,
				paths.PutGameBlob,
				paths.GameParams{Slug: slug},
				map[string]string{"Content-Type": contentType, "Content-Digest": digest},
				nil,
				progressReader,
			)
			if err != nil {
				uploadErrChan <- fmt.Errorf("failed API request: %w", err)
				return
			}
			defer resp.Body.Close()

			uploadErrChan <- nil
		}()

		var zipErr, uploadErr error
		select {
		case zipErr = <-zipErrChan:
		case uploadErr = <-uploadErrChan:
		}

		pipeReader.Close()

		if zipErr == nil {
			select {
			case uploadErr = <-uploadErrChan:
			case zipErr = <-zipErrChan:
			}
		} else {
			select {
			case uploadErr = <-uploadErrChan:
			default:
			}
		}

		if zipErr != nil {
			prog.SetError(fmt.Errorf("failed to create zip: %w", zipErr))
			return
		}
		if uploadErr != nil {
			prog.SetError(uploadErr)
			return
		}

		prog.MarkComplete()
	}()

	return prog, nil
}

func zipDirectoryToWriter(directory string, writer io.Writer) error {
	zipWriter := zip.NewWriter(writer)
	defer zipWriter.Close()

	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == directory {
			return nil
		}

		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create zip header: %w", err)
		}
		header.Name = filepath.ToSlash(relPath)

		if info.IsDir() {
			header.Name += "/"
			_, err := zipWriter.CreateHeader(header)
			return err
		}

		fileWriter, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create zip writer: %w", err)
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		_, err = io.Copy(fileWriter, file)
		if err != nil {
			return fmt.Errorf("failed to copy file to zip: %w", err)
		}

		return nil
	})
}
