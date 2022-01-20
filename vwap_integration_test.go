// +build integration

package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCalculateBTCUSDVwap(t *testing.T) {
	tradePair := "BTC-USD"
	vwape, err := NewVwapEngine([]string{tradePair})
	assert.Nil(t, err)

	go vwape.Calculate()
	time.Sleep(5 * time.Second)

	assert.Greater(t, vwape.storageManagers[tradingPair(tradePair)].GetSum(), 0.00)
}
