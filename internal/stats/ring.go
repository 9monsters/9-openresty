// Package stats provides an in-memory ring buffer for recent system and nginx
// metrics, so the dashboard can render time-series charts without hitting the
// database on every poll.
package stats

import (
	"sync"
	"time"
)

// SystemPoint is one system-info snapshot.
type SystemPoint struct {
	Time         time.Time `json:"time"`
	CPUCount     int       `json:"cpuCount"`
	MemoryUsed   uint64    `json:"memoryUsed"`
	MemoryTotal  uint64    `json:"memoryTotal"`
	LoadAvg1     float64   `json:"loadAvg1"`
	LoadAvg5     float64   `json:"loadAvg5"`
	LoadAvg15    float64   `json:"loadAvg15"`
	UptimeSeconds uint64   `json:"uptimeSeconds"`
}

// NginxPoint is one nginx-stats snapshot.
type NginxPoint struct {
	Time              time.Time `json:"time"`
	ActiveConnections int64     `json:"activeConnections"`
	Accepts           int64     `json:"accepts"`
	Handled           int64     `json:"handled"`
	Requests          int64     `json:"requests"`
	Reading           int64     `json:"reading"`
	Writing           int64     `json:"writing"`
	Waiting           int64     `json:"waiting"`
}

// Ring holds a fixed-size circular buffer of snapshots.
type Ring[T any] struct {
	mu     sync.RWMutex
	buf    []T
	head   int
	size   int
	filled bool
}

// NewRing creates a buffer holding up to capacity entries.
func NewRing[T any](capacity int) *Ring[T] {
	return &Ring[T]{buf: make([]T, capacity)}
}

// Push adds a snapshot, overwriting the oldest when full.
func (r *Ring[T]) Push(v T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buf[r.head] = v
	r.head = (r.head + 1) % len(r.buf)
	if r.size < len(r.buf) {
		r.size++
	} else {
		r.filled = true
	}
}

// Snapshot returns a chronological copy of the buffer contents.
func (r *Ring[T]) Snapshot() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.size == 0 {
		return nil
	}
	out := make([]T, r.size)
	if !r.filled {
		copy(out, r.buf[:r.size])
		return out
	}
	start := r.head % len(r.buf)
	n := copy(out, r.buf[start:])
	copy(out[n:], r.buf[:start])
	return out
}

var (
	SystemHistory = NewRing[SystemPoint](600) // 10 min @ 1s
	NginxHistory  = NewRing[NginxPoint](600)
)