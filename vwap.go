package main

import (
	"fmt"
	"strconv"

	sm "github.com/fernandovaladao/vwap/storage_manager"
	ts "github.com/fernandovaladao/vwap/trade_streaming"

	log "github.com/sirupsen/logrus"
)

const maxStorage = 200

type tradingPair string

type VwapEngine struct {
	storageManagers map[tradingPair]sm.StorageManager
	client          ts.Client
}

func NewVwapEngine(tradingPairs []string) (*VwapEngine, error) {
	storageManagers := make(map[tradingPair]sm.StorageManager)
	for _, tp := range tradingPairs {
		storageManagers[tradingPair(tp)] = sm.NewInMemoryStorageManager()
	}
	if streamingClient, err := ts.NewCoinbaseClient(tradingPairs); err != nil {
		return nil, err
	} else {
		return &VwapEngine{
			storageManagers: storageManagers,
			client:          streamingClient,
		}, nil
	}
}

func (vwape *VwapEngine) Calculate() {
	for {
		trade, err := vwape.readNextTradePair()
		if err != nil {
			continue
		}
		vwape.logUpdatedVwap(trade)
	}
}

func (vwape *VwapEngine) readNextTradePair() (*ts.Trade, error) {
	trade, err := vwape.client.ReadValue()
	if err != nil {
		log.WithError(err).Error("Error reading next trade stream message")
		return nil, err
	}
	if trade.Type == "error" {
		log.WithField("reason", trade.Reason).WithField("message", trade.Message).Error("Error message returned by trade stream client")
		return nil, fmt.Errorf("%s", trade.Type)
	}
	return trade, nil
}

func (vwape *VwapEngine) logUpdatedVwap(tp *ts.Trade) {
	pair := tp.Pair
	if sm, ok := vwape.storageManagers[tradingPair(pair)]; !ok {
		if pair != "" {
			log.WithField("trade_pair", pair).Warn("Received a trading price for an unexpected trading pair.")
		} else {
			log.WithField("trade", tp).Warn("Unknown trading message")
		}
	} else if price, err := strconv.ParseFloat(tp.Price, 64); err != nil {
		log.WithError(err).WithField("price", tp.Price).Error("Error converting price to float64")
	} else if vwap, err := vwape.calculateVwap(sm, price); err != nil {
		log.WithError(err).WithField("trade_pair", pair).WithField("price", price).Error("Error calculating vwap")
	} else {
		log.WithFields(log.Fields{
			"trade_pair": pair,
			"vwap":       vwap,
		}).Info()
	}
}

func (vwape *VwapEngine) calculateVwap(sm sm.StorageManager, price float64) (float64, error) {
	if err := sm.Store(price); err != nil {
		return 0.00, err
	} else {
		return sm.GetSum() / maxStorage, nil
	}
}
