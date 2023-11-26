package network

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/seternate/go-lanty/pkg/util"
)

type Download struct {
	httpclient     *http.Client
	url            url.URL
	filename       string
	filesize       uint64
	alreadywritten uint64
	startTime      time.Time
	endTime        time.Time
	Done           chan struct{}
	subscriber     []chan struct{}
	Err            error
	mutex          sync.RWMutex
}

func NewDownload(url url.URL) (download *Download, err error) {
	request, err := http.NewRequest(http.MethodHead, url.String(), nil)
	if err != nil {
		return
	}

	httpclient := &http.Client{}
	response, err := httpclient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if response.StatusCode > http.StatusBadRequest {
		err = errors.New(response.Status)
		return
	}

	filesize := response.ContentLength
	if filesize < 0 {
		err = errors.New("missing filesize for download")
		return
	}

	contentDisposition := response.Header.Get("Content-Disposition")
	if !strings.HasPrefix(contentDisposition, "attachment") {
		err = errors.New("no attachment to be downloaded")
		return
	}

	_, filename, foundFilename := strings.Cut(contentDisposition, "filename=")
	if foundFilename && strings.TrimSpace(filename) != "" {
		filename = strings.TrimSpace(filename)
	} else {
		err = errors.New("missing filename for download")
		return
	}

	download = &Download{
		httpclient: httpclient,
		url:        url,
		filename:   filename,
		filesize:   uint64(filesize),
		Done:       make(chan struct{}),
	}

	return
}

func (download *Download) StartDownload(ctx context.Context, writer io.Writer) {
	go download.Download(ctx, writer)
}

func (download *Download) Download(ctx context.Context, writer io.Writer) {
	download.mutex.RLock()
	client := download.httpclient
	url := download.url.String()
	download.mutex.RUnlock()
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		download.mutex.Lock()
		download.Err = err
		download.endTime = time.Now()
		download.mutex.Unlock()
		close(download.Done)
		download.notifySubscriber()
		return
	}
	response, err := client.Do(request)
	if err != nil {
		download.mutex.Lock()
		download.Err = err
		download.endTime = time.Now()
		download.mutex.Unlock()
		close(download.Done)
		download.notifySubscriber()
		return
	}

	//Buffer of 10MByte
	buffer := make([]byte, 10*1024*1024)
	download.mutex.Lock()
	download.startTime = time.Now()
	download.mutex.Unlock()

	for {
		var read, write int

		if ctx.Err() != nil {
			response.Body.Close()
			download.mutex.Lock()
			download.Err = ctx.Err()
			download.endTime = time.Now()
			download.mutex.Unlock()
			close(download.Done)
			download.notifySubscriber()
			return
		}

		read, errRead := response.Body.Read(buffer)
		if read > 0 {
			write, err = writer.Write(buffer[0:read])
			if err != nil {
				response.Body.Close()
				download.mutex.Lock()
				download.Err = err
				download.endTime = time.Now()
				download.mutex.Unlock()
				close(download.Done)
				download.notifySubscriber()
				return
			}
			if write > 0 {
				download.mutex.Lock()
				download.alreadywritten += uint64(write)
				download.mutex.Unlock()
				download.notifySubscriber()
			}
		}
		if errRead != nil && errRead != io.EOF {
			response.Body.Close()
			download.mutex.Lock()
			download.Err = errRead
			download.endTime = time.Now()
			download.mutex.Unlock()
			close(download.Done)
			download.notifySubscriber()
			return
		} else if errRead != nil && errRead == io.EOF {
			response.Body.Close()
			download.mutex.Lock()
			download.endTime = time.Now()
			download.mutex.Unlock()
			close(download.Done)
			download.notifySubscriber()
			return
		}
	}
}

func (download *Download) Filename() (filename string) {
	download.mutex.RLock()
	filename = download.filename
	download.mutex.RUnlock()
	return
}

func (download *Download) Filesize() (size uint64) {
	download.mutex.RLock()
	size = download.filesize
	download.mutex.RUnlock()
	return
}

func (download *Download) StartTime() (time time.Time) {
	download.mutex.RLock()
	time = download.startTime
	download.mutex.RUnlock()
	return
}

func (download *Download) EndTime() (time time.Time) {
	download.mutex.RLock()
	time = download.endTime
	download.mutex.RUnlock()
	return
}

func (download *Download) IsComplete() bool {
	select {
	case <-download.Done:
		return true
	default:
		return false
	}
}

func (download *Download) Progress() (progress float64) {
	filesize := download.Filesize()
	if filesize == 0 {
		return
	}
	download.mutex.RLock()
	progress = float64(download.alreadywritten) / float64(filesize)
	download.mutex.RUnlock()
	return
}

func (download *Download) Duration() (duration time.Duration) {
	finished := download.IsComplete()
	download.mutex.RLock()
	if finished {
		duration = download.endTime.Sub(download.startTime)
	} else {
		duration = time.Since(download.startTime)
	}
	download.mutex.RUnlock()
	return
}

func (download *Download) BytesPerSecond() (speed float64) {
	duration := download.Duration()
	download.mutex.RLock()
	speed = float64(download.alreadywritten) / duration.Seconds()
	download.mutex.RUnlock()
	return
}

func (download *Download) Subscribe(subscriber chan struct{}) {
	download.mutex.Lock()
	download.subscriber = append(download.subscriber, subscriber)
	download.mutex.Unlock()
}

func (download *Download) Unsubscribe(subscriber chan struct{}) {
	download.mutex.Lock()
	index := slices.Index(download.subscriber, subscriber)
	if index >= 0 {
		slices.Delete(download.subscriber, index, index+1)
	}
	download.mutex.Unlock()
}

func (download *Download) notifySubscriber() {
	download.mutex.RLock()
	for _, subscriber := range download.subscriber {
		util.ChannelWriteNonBlocking(subscriber, struct{}{})
	}
	download.mutex.RUnlock()
}
