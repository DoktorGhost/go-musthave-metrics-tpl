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
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.gauge[nameMetric] = value
}

func (ms *MemStorage) UpdateCounter(nameMetric string, value int64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	_, ok := ms.counter[nameMetric]
	if !ok {
		ms.counter[nameMetric] += value
	} else {
		ms.counter[nameMetric] = value
	}
}

func (ms *MemStorage) Read(nameType, nameMetric string) interface{} {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	if nameType == "gauge" {
		return ms.gauge[nameMetric]
	} else if nameType == "counter" {
		return ms.counter[nameMetric]
	}
	return nil
}

func (ms *MemStorage) ReadAll() map[string]interface{} {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	result := make(map[string]interface{})
	for k, v := range ms.gauge {
		result[k] = v
	}
	for k, v := range ms.counter {
		result[k] = v
	}
	return result
}
