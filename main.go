package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	tradePairs := os.Args[1:]
	vwape, err := NewVwapEngine(tradePairs)
	if err != nil {
		log.Fatalf("Error when initializing vwap engine: %v", err)
	}
	vwape.Calculate()
}
