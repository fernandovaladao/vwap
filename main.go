package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	vwape, err := NewVwapEngine([]string{"BTC-USD", "ETH-USD", "ETH-BTC"})
	if err != nil {
		log.Fatalf("Error when initializing vwap: %v", err)
	}
	vwape.Calculate()
}
