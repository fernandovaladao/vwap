package storage_manager

import (
	ds "github.com/fernandovaladao/vwap/data_structures"
)

type InMemoryStorageManager struct {
	buffer ds.Queue
	sum    float64
}

func NewInMemoryStorageManager() *InMemoryStorageManager {
	return &InMemoryStorageManager{
		buffer: ds.NewCircularQueue(),
	}
}

func (sm *InMemoryStorageManager) Store(price float64) error {
	if sm.buffer.IsFull() {
		if elem, err := sm.buffer.Dequeue(); err != nil {
			return err
		} else {
			sm.sum -= elem
		}
	}
	if err := sm.buffer.Enqueue(price); err != nil {
		return err
	}
	sm.sum += price
	return nil
}

func (sm *InMemoryStorageManager) GetSum() float64 {
	return sm.sum
}
