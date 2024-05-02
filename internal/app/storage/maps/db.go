package maps

import (
	"sync"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
	mu      sync.RWMutex
}

func NewMapStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string]int64),
	}
}

func (ms *MemStorage) UpdateGauage(nameMetric string, value float64) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	ms.gauge[nameMetric] = value
}

func (ms *MemStorage) UpdateCounter(nameMetric string, value int64) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	_, ok := ms.counter[nameMetric]
	if !ok {
		ms.counter[nameMetric] += value
	} else {
		ms.counter[nameMetric] = value
	}
}
