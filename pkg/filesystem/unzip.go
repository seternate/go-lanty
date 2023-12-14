package filesystem

import (
	"archive/zip"
	"os"
	"path/filepath"
	"slices"
	"sync/atomic"
	"time"

	"github.com/seternate/go-lanty/pkg/util"
)

type Unzip struct {
	alreadywritten uint64
	subscriber     []chan float64

	filename    string
	Destination string
	startTime   time.Time
	endTime     time.Time
	filesize    int64
	Err         error
	Done        chan struct{}
}

func NewUnzip(filename string, destination string) (unzip *Unzip) {
	unzip = &Unzip{
		alreadywritten: 0,
		filename:       filename,
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
		unzip.notifySubscriber()
		unzip.Err = err
	}()

	archive, err := zip.OpenReader(unzip.Filename())
	if err != nil {
		close(unzip.Done)
		return
	}

	archiveReaderFile := archive.Reader.File
	alreadywritten := uint64(0)
	unzip.startTime = time.Now()
	buffer := make([]byte, 1024)

	for _, file := range archiveReaderFile {
		reader, err := file.Open()
		if err != nil {
			close(unzip.Done)
			archive.Close()
			return err
		}
		//defer reader.Close()

		path := filepath.Join(unzip.Destination, file.Name)
		_ = os.Remove(path)
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			close(unzip.Done)
			archive.Close()
			reader.Close()
			return err
		}
		if file.FileInfo().IsDir() {
			reader.Close()
			continue
		}
		err = os.Remove(path)
		if err != nil {
			close(unzip.Done)
			archive.Close()
			reader.Close()
			return err
		}

		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			close(unzip.Done)
			archive.Close()
			reader.Close()
			return err
		}
		defer writer.Close()

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
				reader.Close()
				break
			}
		}
		reader.Close()
	}
	unzip.endTime = time.Now()
	close(unzip.Done)
	archive.Close()
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
	if unzip.Filesize() <= 0 {
		return 0
	}
	return float64(unzip.alreadywritten) / float64(unzip.Filesize())
}

func (unzip *Unzip) Duration() time.Duration {
	if unzip.IsComplete() {
		return unzip.EndTime().Sub(unzip.StartTime())
	}

	return time.Since(unzip.StartTime())
}

func (unzip *Unzip) BytesPerSecond() float64 {
	return float64(unzip.alreadywritten) / unzip.Duration().Seconds()
}

func (unzip *Unzip) Filesize() int64 {
	return unzip.filesize
}

func (unzip *Unzip) Filename() string {
	return unzip.filename
}

func (unzip *Unzip) StartTime() time.Time {
	return unzip.startTime
}

func (unzip *Unzip) EndTime() time.Time {
	return unzip.endTime
}

func (unzip *Unzip) Subscribe(subscriber chan float64) {
	unzip.subscriber = append(unzip.subscriber, subscriber)
}

func (unzip *Unzip) Unsubscribe(subscriber chan float64) {
	index := slices.Index(unzip.subscriber, subscriber)
	slices.Delete(unzip.subscriber, index, index+1)
}

func (unzip *Unzip) notifySubscriber() {
	for _, subscriber := range unzip.subscriber {
		util.ChannelWriteNonBlocking(subscriber, unzip.Progress())
	}
}

func (unzip *Unzip) setSize() {
	archive, err := zip.OpenReader(unzip.Filename())
	if err != nil {
		return
	}
	defer archive.Close()

	size := int64(0)

	for _, file := range archive.Reader.File {
		size += int64(file.UncompressedSize64)
	}

	unzip.filesize = size
}
