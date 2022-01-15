package data_structures

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmpty(t *testing.T) {
	testCases := []struct {
		name     string
		cq       CircularQueue
		expected bool
	}{
		{
			name:     "Circular Queue is empty",
			cq:       CircularQueue{},
			expected: true,
		},
		{
			name: "Circular Queue is not empty",
			cq: CircularQueue{
				elements:  [MaxSize]float64{1.0},
				rearIndex: 1,
			},
			expected: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.cq.IsEmpty())
		})
	}
}

func TestIsFull(t *testing.T) {
	testCases := []struct {
		name     string
		cq       CircularQueue
		expected bool
	}{
		{
			name: "Circular Queue is not full",
			cq: CircularQueue{
				elements:  [MaxSize]float64{1.0, 2.0, 3.0},
				rearIndex: 4,
			},
			expected: false,
		},
		{
			name: "Circular Queue is full",
			cq: CircularQueue{
				elements: [MaxSize]float64{ //MaxSize - 1 elements here
				},
				rearIndex: MaxSize - 1,
			},
			expected: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.cq.IsFull())
		})
	}
}

func TestEnqueue(t *testing.T) {
	testCases := []struct {
		name     string
		cq       CircularQueue
		expected error
	}{
		{
			name:     "Enqueue element in non-full queue",
			cq:       CircularQueue{},
			expected: nil,
		},
		{
			name: "Enqueue element in full queue",
			cq: CircularQueue{
				elements: [MaxSize]float64{ //MaxSize - 1 elements here
				},
				rearIndex: MaxSize - 1,
			},
			expected: fmt.Errorf("queue is full"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.cq.Enqueue(1.0))
		})
	}
}

func TestDequeue(t *testing.T) {
	testCases := []struct {
		name     string
		cq       CircularQueue
		expected struct {
			returned float64
			err      error
		}
	}{
		{
			name: "Dequeue element in empty queue",
			cq:   CircularQueue{},
			expected: struct {
				returned float64
				err      error
			}{
				err: fmt.Errorf("queue is empty"),
			},
		},
		{
			name: "Dequeue element in non-empty queue",
			cq: CircularQueue{
				elements:  [MaxSize]float64{1.0, 2.0, 3.0},
				rearIndex: 4,
			},
			expected: struct {
				returned float64
				err      error
			}{
				returned: 1.0,
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			elem, err := tt.cq.Dequeue()
			assert.Equal(t, tt.expected.err, err)
			assert.Equal(t, tt.expected.returned, elem)
		})
	}
}
