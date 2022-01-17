package main

import (
	"fmt"
	"testing"

	sm "github.com/fernandovaladao/vwap/storage_manager"
	ts "github.com/fernandovaladao/vwap/trade_streaming"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestReadNextTradePair(t *testing.T) {
	ctrl := gomock.NewController(t)

	testcases := []struct {
		name       string
		mockClient func(ctrl *gomock.Controller) ts.Client
		doAsserts  func(t *testing.T, trade *ts.Trade, err error)
	}{
		{
			name: "Happy path: no errors returned",
			mockClient: func(ctrl *gomock.Controller) ts.Client {
				client := ts.NewMockClient(ctrl)
				trade := &ts.Trade{
					Type:  "match",
					Pair:  "BTC-USD",
					Price: "10.00",
				}
				client.EXPECT().ReadValue().Return(trade, nil)
				return client
			},
			doAsserts: func(t *testing.T, trade *ts.Trade, err error) {
				expected := &ts.Trade{
					Type:  "match",
					Pair:  "BTC-USD",
					Price: "10.00",
				}
				assert.EqualValues(t, expected, trade)
				assert.Nil(t, err)
			},
		},
		{
			name: "Error msg returned",
			mockClient: func(ctrl *gomock.Controller) ts.Client {
				client := ts.NewMockClient(ctrl)
				trade := &ts.Trade{
					Type:    "error",
					Pair:    "FOO-BAR",
					Reason:  "reason",
					Message: "message",
				}
				client.EXPECT().ReadValue().Return(trade, nil)
				return client
			},
			doAsserts: func(t *testing.T, trade *ts.Trade, err error) {
				assert.Nil(t, trade)
				assert.NotNil(t, err)
			},
		},
		{
			name: "Error reading value",
			mockClient: func(ctrl *gomock.Controller) ts.Client {
				client := ts.NewMockClient(ctrl)
				client.EXPECT().ReadValue().Return(nil, fmt.Errorf("error"))
				return client
			},
			doAsserts: func(t *testing.T, trade *ts.Trade, err error) {
				assert.Nil(t, trade)
				assert.NotNil(t, err)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			vwap := VwapEngine{
				client: tt.mockClient(ctrl),
			}

			trade, err := vwap.readNextTradePair()

			tt.doAsserts(t, trade, err)
		})
	}
}

func TestLogUpdatedVwap(t *testing.T) {
	// TODO check logrus on how to do asserts with logs.
}

func TestCalculateVwap(t *testing.T) {
	ctrl := gomock.NewController(t)

	testcases := []struct {
		name               string
		price              float64
		mockStorageManager func(ctrl *gomock.Controller, price float64) sm.StorageManager
		doAsserts          func(t *testing.T, vwap float64, err error)
	}{
		{
			name: "Happy path: no errors returned",
			mockStorageManager: func(ctrl *gomock.Controller, price float64) sm.StorageManager {
				sm := sm.NewMockStorageManager(ctrl)
				sm.EXPECT().Store(price)
				sm.EXPECT().GetSum().Return(maxStorage * 1000.00)
				return sm
			},
			doAsserts: func(t *testing.T, vwap float64, err error) {
				assert.Equal(t, vwap, 1000.00)
				assert.Nil(t, err)
			},
		},
		{
			name: "Error in storage manager",
			mockStorageManager: func(ctrl *gomock.Controller, price float64) sm.StorageManager {
				sm := sm.NewMockStorageManager(ctrl)
				sm.EXPECT().Store(price).Return(fmt.Errorf("error"))
				return sm
			},
			doAsserts: func(t *testing.T, vwap float64, err error) {
				assert.Equal(t, vwap, 0.00)
				assert.NotNil(t, err)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			price := 1000.00
			sm := tt.mockStorageManager(ctrl, price)
			vwape := VwapEngine{}

			vwap, err := vwape.calculateVwap(sm, price)

			tt.doAsserts(t, vwap, err)
		})
	}
}
