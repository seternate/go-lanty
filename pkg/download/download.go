package download

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

type Download struct {
	Filename    string
	Filesize    int64
	request     http.Request
	Done        chan struct{}
	StartTime   time.Time
	End         time.Time
	writtenSize int64

	sizeUnsafe int64

	CanResume    bool
	DidResume    bool
	ctx          context.Context
	cancel       context.CancelFunc
	fi           os.FileInfo
	optionsKnown bool
	writer       io.Writer
	storeBuffer  bytes.Buffer
	bytesResumed int64
	//transfer     *transfer
	bufferSize int
	err        error
}

func NewDownload(response *http.Response) (*Download, error) {
	d := &Download{}

	contentDisposition := response.Header.Get("Content-Disposition")
	if strings.HasPrefix(contentDisposition, "attachment") == false {
		return nil, errors.New("No attachment to be downloaded")
	}

	_, filename, foundFilename := strings.Cut(contentDisposition, "filename=")
	if foundFilename && strings.TrimSpace(filename) != "" {
		d.Filename = strings.TrimSpace(filename)
	} else {
		d.Filename = uuid.New().String()
	}

	d.Filesize = response.ContentLength
	d.request = *response.Request
	d.request.Method = http.MethodGet
	d.Done = make(chan struct{}, 0)

	return d, nil
}

func (download *Download) Start(writer io.Writer) (err error) {
	httpClient := http.Client{}
	response, _ := httpClient.Do(&download.request)
	defer response.Body.Close()

	buf := make([]byte, 1024)
	written := int64(0)
	download.StartTime = time.Now()
	for {
		nr, er := response.Body.Read(buf)
		if nr > 0 {
			nw, ew := writer.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				atomic.StoreInt64(&download.writtenSize, written)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	download.End = time.Now()

	close(download.Done)

	return nil
}

func (d *Download) IsComplete() bool {
	select {
	case <-d.Done:
		return true
	default:
		return false
	}
}

func (d *Download) Progress() float64 {
	if d.Filesize <= 0 {
		return 0
	}
	return float64(d.writtenSize) / float64(d.Filesize)
}

func (d *Download) Duration() time.Duration {
	if d.IsComplete() {
		return d.End.Sub(d.StartTime)
	}

	return time.Now().Sub(d.StartTime)
}

func (d *Download) BytesPerSecond() float64 {
	return float64(d.writtenSize) / d.Duration().Seconds()
}

func (d *Download) Written() int64 {
	return d.writtenSize
}
