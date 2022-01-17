package storage_manager

import (
	"fmt"
	"testing"

	ds "github.com/fernandovaladao/vwap/data_structures"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	ctrl := gomock.NewController(t)

	testCases := []struct {
		name      string
		sum       float64
		price     float64
		mockQueue func(ctrl *gomock.Controller) ds.Queue
		doAsserts func(t *testing.T, err error, price float64, sm InMemoryStorageManager)
	}{
		{
			name:  "Buffer is empty",
			sum:   0.00,
			price: 100.00,
			mockQueue: func(ctrl *gomock.Controller) ds.Queue {
				queue := ds.NewMockQueue(ctrl)
				queue.EXPECT().IsFull().Return(false)
				price := 100.00
				queue.EXPECT().Enqueue(price).Return(nil)
				return queue
			},
			doAsserts: func(t *testing.T, err error, price float64, sm InMemoryStorageManager) {
				assert.Nil(t, err)
				assert.Equal(t, price, sm.sum)
			},
		},
		{
			name:  "Buffer is not empty neither full",
			sum:   50.00,
			price: 100.00,
			mockQueue: func(ctrl *gomock.Controller) ds.Queue {
				queue := ds.NewMockQueue(ctrl)
				queue.EXPECT().IsFull().Return(false)
				queue.EXPECT().Enqueue(100.00).Return(nil)
				return queue
			},
			doAsserts: func(t *testing.T, err error, price float64, sm InMemoryStorageManager) {
				assert.Nil(t, err)
				assert.Equal(t, 150.00, sm.sum)
			},
		},
		{
			name:  "Buffer is not empty neither full and error is returned by enqueue",
			sum:   50.00,
			price: 100.00,
			mockQueue: func(ctrl *gomock.Controller) ds.Queue {
				queue := ds.NewMockQueue(ctrl)
				queue.EXPECT().IsFull().Return(false)
				queue.EXPECT().Enqueue(100.00).Return(fmt.Errorf("error"))
				return queue
			},
			doAsserts: func(t *testing.T, err error, price float64, sm InMemoryStorageManager) {
				assert.NotNil(t, err)
			},
		},
		{
			name:  "Buffer is full",
			sum:   50.00,
			price: 100.00,
			mockQueue: func(ctrl *gomock.Controller) ds.Queue {
				queue := ds.NewMockQueue(ctrl)
				queue.EXPECT().IsFull().Return(true)
				queue.EXPECT().Dequeue().Return(20.00, nil)
				queue.EXPECT().Enqueue(100.00).Return(nil)
				return queue
			},
			doAsserts: func(t *testing.T, err error, price float64, sm InMemoryStorageManager) {
				assert.Nil(t, err)
				assert.Equal(t, 130.00, sm.sum)
			},
		},
		{
			name:  "Buffer is full and error is returned by dequeue",
			sum:   50.00,
			price: 100.00,
			mockQueue: func(ctrl *gomock.Controller) ds.Queue {
				queue := ds.NewMockQueue(ctrl)
				queue.EXPECT().IsFull().Return(true)
				queue.EXPECT().Dequeue().Return(0.00, fmt.Errorf("error"))
				return queue
			},
			doAsserts: func(t *testing.T, err error, price float64, sm InMemoryStorageManager) {
				assert.NotNil(t, err)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			price := tt.price
			sm := InMemoryStorageManager{
				buffer: tt.mockQueue(ctrl),
				sum:    tt.sum,
			}

			err := sm.Store(price)

			tt.doAsserts(t, err, price, sm)
		})
	}

}

func TestGetSum(t *testing.T) {
	expected := 42.00
	sm := InMemoryStorageManager{
		sum: expected,
	}

	assert.Equal(t, expected, sm.GetSum())
}
