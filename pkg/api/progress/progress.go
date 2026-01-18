package progress

import (
	"io"
	"sync"
	"time"
)

type ProgressUpdate struct {
	BytesDownloaded int64
	TotalBytes      int64
	Percentage      float64
	Speed           float64
	Complete        bool
	Error           error
}

type Progress struct {
	mu              sync.RWMutex
	bytesDownloaded int64
	totalBytes      int64
	startTime       time.Time
	lastUpdateTime  time.Time
	lastBytes       int64
	complete        bool
	err             error
	updateInterval  int64
	updates         chan ProgressUpdate
	once            sync.Once
}

func NewProgress() *Progress {
	now := time.Now()
	return &Progress{
		startTime:      now,
		lastUpdateTime: now,
		totalBytes:     -1,
		updateInterval: 1048576,
		updates:        make(chan ProgressUpdate, 1),
	}
}

func (p *Progress) Updates() <-chan ProgressUpdate {
	return p.updates
}

func (p *Progress) BytesDownloaded() int64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.bytesDownloaded
}

func (p *Progress) TotalBytes() int64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.totalBytes
}

func (p *Progress) percentageLocked() float64 {
	if p.totalBytes == 0 {
		return 100
	}
	if p.totalBytes < 0 {
		return -1
	}
	return float64(p.bytesDownloaded) / float64(p.totalBytes) * 100
}

func (p *Progress) speedLocked() float64 {
	if p.lastUpdateTime.Equal(p.startTime) {
		return 0
	}
	elapsed := time.Since(p.lastUpdateTime).Seconds()
	if elapsed <= 0 {
		return 0
	}
	bytesSinceLastUpdate := p.bytesDownloaded - p.lastBytes
	return float64(bytesSinceLastUpdate) / elapsed
}

func (p *Progress) Percentage() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.percentageLocked()
}

func (p *Progress) Speed() float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.speedLocked()
}

func (p *Progress) IsComplete() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.complete
}

func (p *Progress) Error() error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.err
}

func (p *Progress) sendUpdate() {
	p.mu.RLock()
	update := ProgressUpdate{
		BytesDownloaded: p.bytesDownloaded,
		TotalBytes:      p.totalBytes,
		Percentage:      p.percentageLocked(),
		Speed:           p.speedLocked(),
		Complete:        p.complete,
		Error:           p.err,
	}
	p.mu.RUnlock()

	select {
	case p.updates <- update:
	default:
	}
}

func (p *Progress) SetTotalBytes(total int64) {
	p.mu.Lock()
	p.totalBytes = total
	p.mu.Unlock()
	p.sendUpdate()
}

func (p *Progress) addBytes(n int64) {
	p.mu.Lock()
	p.bytesDownloaded += n

	now := time.Now()
	shouldUpdate := p.bytesDownloaded-p.lastBytes >= p.updateInterval || p.lastUpdateTime.IsZero()
	if shouldUpdate {
		p.lastUpdateTime = now
		p.lastBytes = p.bytesDownloaded
	}
	p.mu.Unlock()

	if shouldUpdate {
		p.sendUpdate()
	}
}

func (p *Progress) MarkComplete() {
	p.mu.Lock()
	alreadyComplete := p.complete
	p.complete = true
	if !p.lastUpdateTime.IsZero() {
		elapsed := time.Since(p.startTime).Seconds()
		if elapsed > 0 {
			p.lastUpdateTime = time.Now()
		}
	}
	p.mu.Unlock()

	if !alreadyComplete {
		p.sendUpdate()
	}
	p.once.Do(func() {
		close(p.updates)
	})
}

func (p *Progress) SetError(err error) {
	p.mu.Lock()
	alreadyComplete := p.complete
	p.err = err
	p.complete = true
	p.mu.Unlock()

	if !alreadyComplete {
		p.sendUpdate()
	}
	p.once.Do(func() {
		close(p.updates)
	})
}

type Reader struct {
	reader   io.Reader
	progress *Progress
}

func NewReader(reader io.Reader, progress *Progress) *Reader {
	return &Reader{
		reader:   reader,
		progress: progress,
	}
}

func (pr *Reader) Read(p []byte) (n int, err error) {
	n, err = pr.reader.Read(p)
	if n > 0 {
		pr.progress.addBytes(int64(n))
	}
	return n, err
}
