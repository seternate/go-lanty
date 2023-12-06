package network

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/seternate/go-lanty/pkg/util"
)

type Download struct {
	httpclient     *http.Client
	url            url.URL
	alreadywritten int64
	subscriber     []chan float64

	Filename  string
	Filesize  int64
	Done      chan struct{}
	StartTime time.Time
	EndTime   time.Time
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
		Filename:       filename,
		Filesize:       response.ContentLength,
		Done:           make(chan struct{}),
	}

	return
}

func (download *Download) StartDownload(writer io.Writer) {
	go download.Download(writer)
}

func (download *Download) Download(writer io.Writer) (err error) {
	defer func() {
		close(download.Done)
		download.notifySubscriber()
		download.Err = err
	}()

	request, err := http.NewRequest(http.MethodGet, download.url.String(), nil)
	if err != nil {
		return
	}
	response, err := download.httpclient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	buffer := make([]byte, 1024)
	alreadywritten := int64(0)
	download.StartTime = time.Now()

	for {
		var errRead, errWrite error
		var read, write int

		read, errRead = response.Body.Read(buffer)
		if read > 0 {
			write, errWrite = writer.Write(buffer[0:read])
			if write > 0 {
				alreadywritten += int64(write)
				atomic.StoreInt64(&download.alreadywritten, alreadywritten)
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
	download.EndTime = time.Now()
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
	if download.Filesize <= 0 {
		return 0
	}
	return float64(download.alreadywritten) / float64(download.Filesize)
}

func (download *Download) Duration() time.Duration {
	if download.IsComplete() {
		return download.EndTime.Sub(download.StartTime)
	}

	return time.Since(download.StartTime)
}

func (download *Download) BytesPerSecond() float64 {
	return float64(download.alreadywritten) / download.Duration().Seconds()
}

func (download *Download) Subscribe(subscriber chan float64) {
	download.subscriber = append(download.subscriber, subscriber)
}

func (download *Download) notifySubscriber() {
	for _, subscriber := range download.subscriber {
		util.ChannelWriteNonBlocking(subscriber, download.Progress())
	}
}
