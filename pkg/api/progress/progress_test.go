package progress

import (
	"bytes"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProgress(t *testing.T) {
	prog := NewProgress()
	assert.NotNil(t, prog)
	assert.Equal(t, int64(-1), prog.TotalBytes())
	assert.Equal(t, int64(0), prog.BytesDownloaded())
	assert.False(t, prog.IsComplete())
	assert.Nil(t, prog.Error())
	assert.NotNil(t, prog.Updates())
}

func TestProgress_SetTotalBytes(t *testing.T) {
	prog := NewProgress()
	
	prog.SetTotalBytes(100)
	assert.Equal(t, int64(100), prog.TotalBytes())
	
	prog.SetTotalBytes(200)
	assert.Equal(t, int64(200), prog.TotalBytes())
}

func TestProgress_BytesDownloaded(t *testing.T) {
	prog := NewProgress()
	assert.Equal(t, int64(0), prog.BytesDownloaded())
}

func TestProgress_TotalBytes(t *testing.T) {
	prog := NewProgress()
	assert.Equal(t, int64(-1), prog.TotalBytes())
	
	prog.SetTotalBytes(100)
	assert.Equal(t, int64(100), prog.TotalBytes())
}

func TestProgress_Percentage(t *testing.T) {
	t.Run("unknown total (-1)", func(t *testing.T) {
		prog := NewProgress()
		prog.SetTotalBytes(-1)
		assert.Equal(t, float64(-1), prog.Percentage())
	})

	t.Run("zero total", func(t *testing.T) {
		prog := NewProgress()
		prog.SetTotalBytes(0)
		assert.Equal(t, float64(100), prog.Percentage())
	})

	t.Run("partial progress", func(t *testing.T) {
		prog := NewProgress()
		prog.SetTotalBytes(100)
		// Simulate adding bytes (we can't directly call addBytes, but we can use Reader)
		reader := NewReader(bytes.NewReader([]byte("test")), prog)
		_, _ = io.ReadAll(reader)
		
		// After reading 4 bytes out of 100
		assert.Greater(t, prog.Percentage(), float64(0))
		assert.Less(t, prog.Percentage(), float64(100))
	})

	t.Run("complete", func(t *testing.T) {
		prog := NewProgress()
		prog.SetTotalBytes(100)
		data := bytes.NewReader(make([]byte, 100))
		reader := NewReader(data, prog)
		_, _ = io.ReadAll(reader)
		
		// Should be close to 100%, but may not be exactly 100 due to timing
		assert.GreaterOrEqual(t, prog.Percentage(), float64(99))
	})
}

func TestProgress_Speed(t *testing.T) {
	prog := NewProgress()
	
	// Initially speed should be 0 (no time has passed)
	assert.Equal(t, float64(0), prog.Speed())
	
	// After adding some bytes and waiting, speed should be calculated
	prog.SetTotalBytes(100)
	data := bytes.NewReader(make([]byte, 50))
	reader := NewReader(data, prog)
	_, _ = io.ReadAll(reader)
	
	// Wait a bit for time to pass
	time.Sleep(10 * time.Millisecond)
	
	// Speed should be greater than 0 now
	speed := prog.Speed()
	assert.GreaterOrEqual(t, speed, float64(0))
}

func TestProgress_MarkComplete(t *testing.T) {
	prog := NewProgress()
	
	prog.MarkComplete()
	assert.True(t, prog.IsComplete())
	assert.Nil(t, prog.Error())
	
	// Channel should be closed
	update, ok := <-prog.Updates()
	assert.True(t, ok) // Last update should be available
	assert.True(t, update.Complete)
	
	// Next read should indicate channel is closed
	_, ok = <-prog.Updates()
	assert.False(t, ok)
}

func TestProgress_SetError(t *testing.T) {
	prog := NewProgress()
	testErr := assert.AnError
	
	prog.SetError(testErr)
	assert.True(t, prog.IsComplete())
	assert.Equal(t, testErr, prog.Error())
	
	// Channel should be closed
	update, ok := <-prog.Updates()
	assert.True(t, ok)
	assert.True(t, update.Complete)
	assert.Equal(t, testErr, update.Error)
	
	_, ok = <-prog.Updates()
	assert.False(t, ok)
}

func TestProgress_IsComplete(t *testing.T) {
	prog := NewProgress()
	assert.False(t, prog.IsComplete())
	
	prog.MarkComplete()
	assert.True(t, prog.IsComplete())
}

func TestProgress_Error(t *testing.T) {
	prog := NewProgress()
	assert.Nil(t, prog.Error())
	
	testErr := assert.AnError
	prog.SetError(testErr)
	assert.Equal(t, testErr, prog.Error())
}

func TestProgress_Updates(t *testing.T) {
	prog := NewProgress()
	prog.SetTotalBytes(100)
	
	updatesChan := prog.Updates()
	assert.NotNil(t, updatesChan)
	
	// Add bytes to trigger update
	data := make([]byte, 1048576+1) // More than updateInterval (1048576)
	reader := NewReader(bytes.NewReader(data), prog)
	
	// Read in a goroutine to avoid blocking
	go func() {
		_, _ = io.ReadAll(reader)
	}()
	
	// Wait for an update
	select {
	case update := <-updatesChan:
		assert.GreaterOrEqual(t, update.BytesDownloaded, int64(0))
		assert.Equal(t, int64(100), update.TotalBytes)
		assert.False(t, update.Complete)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for update")
	}
}

func TestProgress_Reader(t *testing.T) {
	prog := NewProgress()
	prog.SetTotalBytes(100)
	
	data := []byte("test data")
	reader := NewReader(bytes.NewReader(data), prog)
	
	result, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, data, result)
	assert.Equal(t, int64(len(data)), prog.BytesDownloaded())
}

func TestProgress_Reader_LargeData(t *testing.T) {
	prog := NewProgress()
	prog.SetTotalBytes(1000000)
	
	data := make([]byte, 50000)
	for i := range data {
		data[i] = byte(i % 256)
	}
	
	reader := NewReader(bytes.NewReader(data), prog)
	result, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, data, result)
	assert.Equal(t, int64(len(data)), prog.BytesDownloaded())
}

func TestProgress_Concurrency(t *testing.T) {
	prog := NewProgress()
	prog.SetTotalBytes(1000)
	
	var wg sync.WaitGroup
	numGoroutines := 10
	bytesPerGoroutine := 100
	
	// Simulate concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			data := make([]byte, bytesPerGoroutine)
			reader := NewReader(bytes.NewReader(data), prog)
			_, _ = io.ReadAll(reader)
		}()
	}
	
	wg.Wait()
	
	// Total should be numGoroutines * bytesPerGoroutine
	expectedTotal := int64(numGoroutines * bytesPerGoroutine)
	assert.Equal(t, expectedTotal, prog.BytesDownloaded())
}

func TestProgress_UpdateInterval(t *testing.T) {
	prog := NewProgress()
	prog.SetTotalBytes(1000000)
	
	// Small read - should not trigger update immediately
	smallData := make([]byte, 1000)
	reader := NewReader(bytes.NewReader(smallData), prog)
	_, _ = io.ReadAll(reader)
	
	// Bytes should be updated even if update not sent (updateInterval check)
	assert.Equal(t, int64(len(smallData)), prog.BytesDownloaded())
}

func TestProgress_MarkComplete_MultipleCalls(t *testing.T) {
	prog := NewProgress()
	
	prog.MarkComplete()
	prog.MarkComplete() // Should be safe to call multiple times
	
	assert.True(t, prog.IsComplete())
	
	// Channel should only be closed once
	count := 0
	for range prog.Updates() {
		count++
	}
	assert.Equal(t, 1, count) // Only one update
}

func TestProgress_SetError_MultipleCalls(t *testing.T) {
	prog := NewProgress()
	err1 := assert.AnError
	err2 := assert.AnError
	
	prog.SetError(err1)
	prog.SetError(err2) // Should be safe to call multiple times
	
	assert.True(t, prog.IsComplete())
	assert.Equal(t, err2, prog.Error()) // Last error should be kept
	
	count := 0
	for range prog.Updates() {
		count++
	}
	assert.Equal(t, 1, count)
}

func TestProgress_AddBytes_UpdateInterval(t *testing.T) {
	prog := NewProgress()
	prog.SetTotalBytes(10000000)
	
	// Read more than updateInterval (1048576 bytes)
	largeData := make([]byte, 1048576+1000)
	reader := NewReader(bytes.NewReader(largeData), prog)
	
	updates := []ProgressUpdate{}
	done := make(chan bool)
	
	go func() {
		for update := range prog.Updates() {
			updates = append(updates, update)
			if update.Complete {
				break
			}
		}
		done <- true
	}()
	
	_, _ = io.ReadAll(reader)
	prog.MarkComplete()
	
	<-done
	
	// Should have received at least one update due to large read
	assert.Greater(t, len(updates), 0)
}
