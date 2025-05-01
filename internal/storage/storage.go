package storage

import "sync"

type Storage interface {
	SetGauge(name string, value float64)
	GetGauges() map[string]float64

	SetCounter(name string, value int64)
	GetCounters() map[string]int64
}

type MemStorage struct {
	mu       sync.RWMutex
	gauges   map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (s *MemStorage) SetGauge(name string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gauges[name] = value
}

func (s *MemStorage) GetGauges() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]float64)
	for k, v := range s.gauges {
		copy[k] = v
	}
	return copy
}

func (s *MemStorage) SetCounter(name string, value int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counters[name] = value
}

func (s *MemStorage) GetCounters() map[string]int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	copy := make(map[string]int64)
	for k, v := range s.counters {
		copy[k] = v
	}
	return copy
}
