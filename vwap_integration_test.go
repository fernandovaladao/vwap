
package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBTCUSDCalculateVwap(t *testing.T) {
	vwape, err := NewVwapEngine([]string{"BTC-USD"})
	assert.Nil(t, err)

	go vwape.Calculate()
	time.Sleep(5*time.Second)

	assert.Greater(t, vwape.storageManagers["BTC-USD"].GetSum(), 0.00)
}
