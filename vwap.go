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
		// TODO once the circular queue used by StorageManager is changed to be a slice instead of a fixed-size array, we should pass "maxStorage"
		// as a parameter here and use it to build the circular queue.
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
		log.WithError(err).Error()
		return nil, err
	}
	if trade.Type == "error" {
		log.WithFields(log.Fields{
			"reason":  trade.Reason,
			"message": trade.Message,
		}).Error()
		return nil, fmt.Errorf("%s", trade.Type)
	}
	return trade, nil
}

func (vwape *VwapEngine) logUpdatedVwap(tp *ts.Trade) {
	pair := tp.Pair
	if sm, ok := vwape.storageManagers[tradingPair(pair)]; !ok {
		log.WithFields(log.Fields{
			"trade_pair": pair,
		}).Warn("Received a trading price for an unexpected trading pair.")
	} else if price, err := strconv.ParseFloat(tp.Price, 64); err != nil {
		log.WithError(err).Errorf("Error converting %s to float64.", tp.Price)
	} else {
		vwap := vwape.calculateVwap(sm, price)
		log.WithFields(log.Fields{
			"trade_pair": pair,
			"vwap":       vwap,
		}).Info()
	}
}

func (vwape *VwapEngine) calculateVwap(sm sm.StorageManager, price float64) float64 {
	sm.Store(price)
	return sm.GetSum() / maxStorage
}
