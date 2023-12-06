package filesystem

import (
	"archive/zip"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/seternate/go-lanty/pkg/util"
)

type Unzip struct {
	alreadywritten uint64
	subscriber     []chan float64

	Filename    string
	Destination string
	Size        uint64
	Err         error
	Done        chan struct{}
}

func NewUnzip(filename string, destination string) (unzip *Unzip) {
	unzip = &Unzip{
		alreadywritten: 0,
		Filename:       filename,
		Destination:    destination,
		Done:           make(chan struct{}),
	}
	unzip.setSize()

	return
}

func (unzip *Unzip) StartUnzip() {
	go unzip.Unzip()
}

func (unzip *Unzip) Unzip() (err error) {
	defer func() {
		unzip.Err = err
		close(unzip.Done)
		unzip.notifySubscriber()
	}()

	archive, err := zip.OpenReader(unzip.Filename)
	if err != nil {
		return
	}
	defer archive.Close()

	archiveReaderFile := archive.Reader.File
	alreadywritten := uint64(0)

	for _, file := range archiveReaderFile {
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()

		path := filepath.Join(unzip.Destination, file.Name)
		_ = os.Remove(path)
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		if file.FileInfo().IsDir() {
			continue
		}
		err = os.Remove(path)
		if err != nil {
			return err
		}

		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer writer.Close()

		buffer := make([]byte, 1024)
		for {
			var read, write int

			read, err = reader.Read(buffer)
			if read > 0 {
				write, err = writer.Write(buffer[0:read])
				if write > 0 {
					alreadywritten += uint64(write)
					atomic.StoreUint64(&unzip.alreadywritten, alreadywritten)
					unzip.notifySubscriber()
				}
			}
			if err != nil || read != write {
				break
			}
		}
	}

	return
}

func (unzip *Unzip) IsComplete() bool {
	select {
	case <-unzip.Done:
		return true
	default:
		return false
	}
}

func (unzip *Unzip) Progress() float64 {
	if unzip.Size <= 0 {
		return 0
	}
	return float64(unzip.alreadywritten) / float64(unzip.Size)
}

func (unzip *Unzip) Subscribe(subscriber chan float64) {
	unzip.subscriber = append(unzip.subscriber, subscriber)
}

func (unzip *Unzip) notifySubscriber() {
	for _, subscriber := range unzip.subscriber {
		util.ChannelWriteNonBlocking(subscriber, unzip.Progress())
	}
}

func (unzip *Unzip) setSize() {
	archive, err := zip.OpenReader(unzip.Filename)
	if err != nil {
		return
	}
	defer archive.Close()

	size := uint64(0)

	for _, file := range archive.Reader.File {
		size += file.UncompressedSize64
	}

	unzip.Size = size
}
