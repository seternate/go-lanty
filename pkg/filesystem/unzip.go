package filesystem

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/seternate/go-lanty/pkg/util"
)

type Unzip struct {
	filename       string
	filesize       uint64
	Destination    string
	alreadywritten uint64
	startTime      time.Time
	endTime        time.Time
	Done           chan struct{}
	subscriber     []chan struct{}
	Err            error
	mutex          sync.RWMutex
}

func NewUnzip(filename string, destination string) (unzip *Unzip) {
	unzip = &Unzip{
		filename:    filename,
		filesize:    ZipExtractedSize(filename),
		Destination: destination,
		Done:        make(chan struct{}),
	}
	return
}

func (unzip *Unzip) StartUnzip(ctx context.Context) {
	go unzip.Unzip(ctx)
}

func (unzip *Unzip) Unzip(ctx context.Context) {
	if ctx.Err() != nil {
		unzip.mutex.Lock()
		unzip.Err = ctx.Err()
		unzip.endTime = time.Now()
		unzip.mutex.Unlock()
		close(unzip.Done)
		unzip.notifySubscriber()
		return
	}

	archive, err := zip.OpenReader(unzip.Filename())
	if err != nil {
		unzip.mutex.Lock()
		unzip.Err = err
		unzip.endTime = time.Now()
		unzip.mutex.Unlock()
		close(unzip.Done)
		unzip.notifySubscriber()
		return
	}

	archiveReaderFile := archive.Reader.File
	//Buffer of 10MByte
	buffer := make([]byte, 10*1024*1024)
	unzip.mutex.Lock()
	unzip.startTime = time.Now()
	unzip.mutex.Unlock()

	for _, file := range archiveReaderFile {
		reader, err := file.Open()
		if err != nil {
			archive.Close()
			unzip.mutex.Lock()
			unzip.Err = err
			unzip.endTime = time.Now()
			unzip.mutex.Unlock()
			close(unzip.Done)
			unzip.notifySubscriber()
			return
		}

		unzip.mutex.RLock()
		path := filepath.Join(unzip.Destination, file.Name)
		unzip.mutex.RUnlock()
		_ = os.Remove(path)
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			reader.Close()
			archive.Close()
			unzip.mutex.Lock()
			unzip.Err = err
			unzip.endTime = time.Now()
			unzip.mutex.Unlock()
			close(unzip.Done)
			unzip.notifySubscriber()
			return
		}
		if file.FileInfo().IsDir() {
			reader.Close()
			continue
		}
		err = os.Remove(path)
		if err != nil {
			reader.Close()
			archive.Close()
			unzip.mutex.Lock()
			unzip.Err = err
			unzip.endTime = time.Now()
			unzip.mutex.Unlock()
			close(unzip.Done)
			unzip.notifySubscriber()
			return
		}

		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			reader.Close()
			archive.Close()
			unzip.mutex.Lock()
			unzip.Err = err
			unzip.endTime = time.Now()
			unzip.mutex.Unlock()
			close(unzip.Done)
			unzip.notifySubscriber()
			return
		}

		for {
			var read, write int

			if ctx.Err() != nil {
				writer.Close()
				reader.Close()
				archive.Close()
				unzip.mutex.Lock()
				unzip.Err = ctx.Err()
				unzip.endTime = time.Now()
				unzip.mutex.Unlock()
				close(unzip.Done)
				unzip.notifySubscriber()
				return
			}

			read, errRead := reader.Read(buffer)
			if read > 0 {
				write, err = writer.Write(buffer[0:read])
				if err != nil {
					writer.Close()
					reader.Close()
					archive.Close()
					unzip.mutex.Lock()
					unzip.Err = err
					unzip.endTime = time.Now()
					unzip.mutex.Unlock()
					close(unzip.Done)
					unzip.notifySubscriber()
					return
				}
				if write > 0 {
					unzip.mutex.Lock()
					unzip.alreadywritten += uint64(write)
					unzip.mutex.Unlock()
					unzip.notifySubscriber()
				}
			}
			if errRead != nil && errRead != io.EOF {
				writer.Close()
				reader.Close()
				archive.Close()
				unzip.mutex.Lock()
				unzip.Err = err
				unzip.endTime = time.Now()
				unzip.mutex.Unlock()
				close(unzip.Done)
				unzip.notifySubscriber()
				return
			} else if errRead != nil && errRead == io.EOF {
				break
			}
		}
		writer.Close()
		reader.Close()
	}

	archive.Close()
	unzip.mutex.Lock()
	unzip.endTime = time.Now()
	unzip.mutex.Unlock()
	close(unzip.Done)
	unzip.notifySubscriber()
}

func (unzip *Unzip) Filename() (filename string) {
	unzip.mutex.RLock()
	filename = unzip.filename
	unzip.mutex.RUnlock()
	return
}

func (unzip *Unzip) Filesize() (size uint64) {
	unzip.mutex.RLock()
	size = unzip.filesize
	unzip.mutex.RUnlock()
	return
}

func (unzip *Unzip) StartTime() (time time.Time) {
	unzip.mutex.RLock()
	time = unzip.startTime
	unzip.mutex.RUnlock()
	return
}

func (unzip *Unzip) EndTime() (time time.Time) {
	unzip.mutex.RLock()
	time = unzip.endTime
	unzip.mutex.RUnlock()
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

func (unzip *Unzip) Progress() (progress float64) {
	if unzip.Filesize() == 0 {
		return
	}
	unzip.mutex.RLock()
	progress = float64(unzip.alreadywritten) / float64(unzip.filesize)
	unzip.mutex.RUnlock()
	return
}

func (unzip *Unzip) Duration() (duration time.Duration) {
	finished := unzip.IsComplete()
	unzip.mutex.RLock()
	if finished {
		duration = unzip.endTime.Sub(unzip.startTime)
	} else {
		duration = time.Since(unzip.startTime)
	}
	unzip.mutex.RUnlock()
	return
}

func (unzip *Unzip) BytesPerSecond() (speed float64) {
	duration := unzip.Duration()
	unzip.mutex.RLock()
	speed = float64(unzip.alreadywritten) / duration.Seconds()
	unzip.mutex.RUnlock()
	return
}

func (unzip *Unzip) Subscribe(subscriber chan struct{}) {
	unzip.mutex.Lock()
	unzip.subscriber = append(unzip.subscriber, subscriber)
	unzip.mutex.Unlock()
}

func (unzip *Unzip) Unsubscribe(subscriber chan struct{}) {
	unzip.mutex.Lock()
	index := slices.Index(unzip.subscriber, subscriber)
	if index >= 0 {
		slices.Delete(unzip.subscriber, index, index+1)
	}
	unzip.mutex.Unlock()
}

func (unzip *Unzip) notifySubscriber() {
	unzip.mutex.RLock()
	for _, subscriber := range unzip.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
	unzip.mutex.RUnlock()
}

func ZipExtractedSize(path string) (size uint64) {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return
	}
	defer archive.Close()

	for _, file := range archive.Reader.File {
		size += file.UncompressedSize64
	}
	return
}
