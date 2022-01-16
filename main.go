package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	vwape, err := NewVwapEngine([]string{"BTC-USD", "ETH-USD", "ETH-BTC"})
	if err != nil {
		log.Fatalf("Error when initializing vwap: %v", err)
	}
	vwape.Calculate()
}
