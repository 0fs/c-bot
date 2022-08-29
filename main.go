package main

import (
	"log"
)

var currencyGraph CurrencyGraph
var fees map[string]Fee
var cycleDepth = 3
var asset = "USDT"

//var done chan struct{}

func main() {
	initConfig()
	initSpotConnection()

	fees = make(map[string]Fee)
	symbols = make(map[string]*SymbolInfo)
	initFeesMap()
	initCurrencyGraph()

	log.Println("Using asset:", asset)
	log.Printf("Searching cycles with %d deals...", cycleDepth)
	findCycles(currencyGraph.nodes[asset])
	log.Printf("%d cycles has been found\n\n", len(cycles))

	done := make(chan struct{})
	//wsDepth(done, "BTCUSDT")
	log.Println("Listening wsBookTicker...")
	go wsBookTicker(done, asset)

	<-done
}
