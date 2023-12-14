package network

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/seternate/go-lanty/pkg/util"
)

var _ util.Publisher[float64] = &Download{}

type Download struct {
	httpclient     *http.Client
	url            url.URL
	alreadywritten uint64
	subscriber     []chan float64

	filename  string
	filesize  int64
	Done      chan struct{}
	startTime time.Time
	endTime   time.Time
	Err       error
}

func NewDownload(url url.URL) (download *Download, err error) {
	request, err := http.NewRequest(http.MethodHead, url.String(), nil)
	if err != nil {
		return
	}

	httpclient := &http.Client{}

	response, err := httpclient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

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
		httpclient:     httpclient,
		url:            url,
		alreadywritten: 0,
		filename:       filename,
		filesize:       response.ContentLength,
		Done:           make(chan struct{}),
		startTime:      time.Time{},
		endTime:        time.Time{},
		subscriber:     make([]chan float64, 0, 10),
	}

	return
}

func (download *Download) StartDownload(writer io.Writer) {
	go download.Download(writer)
}

func (download *Download) Download(writer io.Writer) (err error) {
	defer func() {
		download.notifySubscriber()
		download.Err = err
	}()

	request, err := http.NewRequest(http.MethodGet, download.url.String(), nil)
	if err != nil {
		close(download.Done)
		return
	}
	response, err := download.httpclient.Do(request)
	if err != nil {
		close(download.Done)
		return
	}
	defer response.Body.Close()

	buffer := make([]byte, 1024)
	alreadywritten := uint64(0)
	download.startTime = time.Now()

	for {
		var errRead, errWrite error
		var read, write int

		read, errRead = response.Body.Read(buffer)
		if read > 0 {
			write, errWrite = writer.Write(buffer[0:read])
			if write > 0 {
				alreadywritten += uint64(write)
				atomic.StoreUint64(&download.alreadywritten, alreadywritten)
				download.notifySubscriber()
			}
		}
		if errRead != nil || errWrite != nil {
			if errRead != io.EOF {
				err = errRead
			}
			break
		} else if read != write {
			err = io.ErrShortWrite
			break
		}
	}
	download.endTime = time.Now()
	close(download.Done)
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

func (download *Download) Progress() float64 {
	if download.Filesize() <= 0 {
		return 0
	}
	return float64(download.alreadywritten) / float64(download.Filesize())
}

func (download *Download) Duration() time.Duration {
	if download.IsComplete() {
		return download.EndTime().Sub(download.StartTime())
	}

	return time.Since(download.StartTime())
}

func (download *Download) BytesPerSecond() float64 {
	return float64(download.alreadywritten) / download.Duration().Seconds()
}

func (download *Download) Filesize() int64 {
	return download.filesize
}

func (download *Download) Filename() string {
	return download.filename
}

func (download *Download) StartTime() time.Time {
	return download.startTime
}

func (download *Download) EndTime() time.Time {
	return download.endTime
}

func (download *Download) Subscribe(subscriber chan float64) {
	download.subscriber = append(download.subscriber, subscriber)
}

func (download *Download) Unsubscribe(subscriber chan float64) {
	index := slices.Index(download.subscriber, subscriber)
	if index >= 0 {
		slices.Delete(download.subscriber, index, index+1)
	}
}

func (download *Download) notifySubscriber() {
	for _, subscriber := range download.subscriber {
		util.ChannelWriteNonBlocking(subscriber, download.Progress())
	}
}
