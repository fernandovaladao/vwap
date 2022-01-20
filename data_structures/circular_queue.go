package data_structures

import "fmt"

const MaxSize int = 201

type CircularQueue struct {
	elements   [MaxSize]float64
	frontIndex int
	rearIndex  int
}

func NewCircularQueue() *CircularQueue {
	return &CircularQueue{}
}

func (cq *CircularQueue) Enqueue(elem float64) error {
	if cq.IsFull() {
		return fmt.Errorf("queue is full")
	}
	cq.elements[cq.rearIndex] = elem
	cq.rearIndex++
	if cq.rearIndex == MaxSize {
		cq.rearIndex = 0
	}
	return nil
}

func (cq *CircularQueue) Dequeue() (float64, error) {
	if cq.IsEmpty() {
		return 0, fmt.Errorf("queue is empty")
	}
	elem := cq.elements[cq.frontIndex]
	cq.frontIndex++
	if cq.frontIndex == MaxSize {
		cq.frontIndex = 0
	}
	return elem, nil
}

func (cq *CircularQueue) IsFull() bool {
	return ((cq.frontIndex - 1) == cq.rearIndex) || ((cq.frontIndex - 1) == (cq.rearIndex - MaxSize))
}

func (cq *CircularQueue) IsEmpty() bool {
	return cq.frontIndex == cq.rearIndex
}
