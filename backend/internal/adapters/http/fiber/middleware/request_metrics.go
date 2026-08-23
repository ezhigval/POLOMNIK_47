package middleware

import (
	"sync"
	"sync/atomic"
	"time"
)

type RequestMetrics struct {
	mu        sync.Mutex
	last      time.Duration
	sum       time.Duration
	count     uint64
	status5xx uint64
	samples   []time.Duration
}

func NewRequestMetrics() *RequestMetrics {
	return &RequestMetrics{samples: make([]time.Duration, 0, 64)}
}

func (m *RequestMetrics) Observe(d time.Duration) {
	m.ObserveStatus(d, 0)
}

func (m *RequestMetrics) ObserveStatus(d time.Duration, status int) {
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.last = d
	m.sum += d
	atomic.AddUint64(&m.count, 1)
	if status >= 500 && status <= 599 {
		atomic.AddUint64(&m.status5xx, 1)
	}
	if len(m.samples) == 64 {
		m.samples = m.samples[1:]
	}
	m.samples = append(m.samples, d)
}

func (m *RequestMetrics) Status5xx() uint64 {
	if m == nil {
		return 0
	}
	return atomic.LoadUint64(&m.status5xx)
}

func (m *RequestMetrics) Snapshot() (last time.Duration, avg time.Duration, count uint64) {
	if m == nil {
		return 0, 0, 0
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	count = atomic.LoadUint64(&m.count)
	if count == 0 {
		return 0, 0, 0
	}
	return m.last, m.sum / time.Duration(count), count
}
