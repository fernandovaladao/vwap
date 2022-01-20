package main

import (
	"fmt"
	"strconv"
	"testing"

	sm "github.com/fernandovaladao/vwap/storage_manager"
	ts "github.com/fernandovaladao/vwap/trade_streaming"

	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
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
				log:    log.New(),
			}

			trade, err := vwap.readNextTradePair()

			tt.doAsserts(t, trade, err)
		})
	}
}

func TestLogUpdatedVwap(t *testing.T) {
	ctrl := gomock.NewController(t)

	testcases := []struct {
		name               string
		tp                 *ts.Trade
		mockStorageManager func(ctrl *gomock.Controller, tp *ts.Trade) map[tradingPair]sm.StorageManager
		doAsserts          func(t *testing.T, tp *ts.Trade, hook *logTest.Hook)
	}{
		{
			name: "Happy path: no errors returned",
			tp: &ts.Trade{
				Pair: "BTC-USD",
				Price: "100.00",
			},
			mockStorageManager: func(ctrl *gomock.Controller, tp *ts.Trade) map[tradingPair]sm.StorageManager {
				sms := make(map[tradingPair]sm.StorageManager)
				sm := sm.NewMockStorageManager(ctrl)
				price, _ := strconv.ParseFloat(tp.Price, 64)
				sm.EXPECT().Store(price)
				sm.EXPECT().GetSum().Return(float64(price * 1000.00))
				sms[tradingPair(tp.Pair)] = sm
				return sms
			},
			doAsserts: func(t *testing.T, tp *ts.Trade, hook *logTest.Hook){
				assert.Equal(t, 1, len(hook.Entries))
				entry := hook.LastEntry()
				assert.Equal(t, log.InfoLevel, entry.Level)
				assert.Equal(t, tp.Pair, entry.Data["trade_pair"])
				assert.Equal(t, 500.00, entry.Data["vwap"])
			},
		},
		{
			name: "Unexpected trading pair",
			tp: &ts.Trade{
				Pair: "FOO-BAR",
			},
			mockStorageManager: func(ctrl *gomock.Controller, tp *ts.Trade) map[tradingPair]sm.StorageManager {
				sms := make(map[tradingPair]sm.StorageManager)
				sm := sm.NewMockStorageManager(ctrl)
				sms["BTC-USD"] = sm
				return sms
			},
			doAsserts: func(t *testing.T, tp *ts.Trade, hook *logTest.Hook){
				assert.Equal(t, 1, len(hook.Entries))
				entry := hook.LastEntry()
				assert.Equal(t, log.WarnLevel, entry.Level)
				assert.Equal(t, tp.Pair, entry.Data["trade_pair"])
				assert.NotEmpty(t, entry.Message)
			},
		},
		{
			name: "Empty trading pair",
			tp: &ts.Trade{
				Pair: "",
			},
			mockStorageManager: func(ctrl *gomock.Controller, tp *ts.Trade) map[tradingPair]sm.StorageManager {
				sms := make(map[tradingPair]sm.StorageManager)
				sm := sm.NewMockStorageManager(ctrl)
				sms["BTC-USD"] = sm
				return sms
			},
			doAsserts: func(t *testing.T, tp *ts.Trade, hook *logTest.Hook){
				assert.Equal(t, 1, len(hook.Entries))
				entry := hook.LastEntry()
				assert.Equal(t, log.WarnLevel, entry.Level)
				assert.Equal(t, tp, entry.Data["trade"])
				assert.NotEmpty(t, entry.Message)
			},
		},	
		{
			name: "Trading pair with invalid price",
			tp: &ts.Trade{
				Pair: "BTC-USD",
				Price: "XYZ",
			},
			mockStorageManager: func(ctrl *gomock.Controller, tp *ts.Trade) map[tradingPair]sm.StorageManager {
				sms := make(map[tradingPair]sm.StorageManager)
				sm := sm.NewMockStorageManager(ctrl)
				sms[tradingPair(tp.Pair)] = sm
				return sms
			},
			doAsserts: func(t *testing.T, tp *ts.Trade, hook *logTest.Hook){
				assert.Equal(t, 1, len(hook.Entries))
				entry := hook.LastEntry()
				assert.Equal(t, log.ErrorLevel, entry.Level)
				assert.Equal(t, tp.Price, entry.Data["price"])
				assert.NotEmpty(t, entry.Message)
			},
		},
		{
			name: "Error to store in storage manager",
			tp: &ts.Trade{
				Pair: "BTC-USD",
				Price: "100.00",
			},
			mockStorageManager: func(ctrl *gomock.Controller, tp *ts.Trade) map[tradingPair]sm.StorageManager {
				sms := make(map[tradingPair]sm.StorageManager)
				sm := sm.NewMockStorageManager(ctrl)
				price, _ := strconv.ParseFloat(tp.Price, 64)
				sm.EXPECT().Store(price).Return(fmt.Errorf("error"))
				sms[tradingPair(tp.Pair)] = sm
				return sms
			},
			doAsserts: func(t *testing.T, tp *ts.Trade, hook *logTest.Hook){
				assert.Equal(t, 1, len(hook.Entries))
				entry := hook.LastEntry()
				assert.Equal(t, log.ErrorLevel, entry.Level)
				assert.Equal(t, tp.Pair, entry.Data["trade_pair"])
				assert.Equal(t, 100.00, entry.Data["price"])
				assert.NotEmpty(t, entry.Message)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			logTest, hook := logTest.NewNullLogger()
			vwape := VwapEngine{
				log: logTest,
				storageManagers: tt.mockStorageManager(ctrl, tt.tp),
			}

			vwape.logUpdatedVwap(tt.tp)

			tt.doAsserts(t, tt.tp, hook)
		})
	}
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
