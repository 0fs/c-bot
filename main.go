package main

import (
	"log"
)

var currencyGraph CurrencyGraph
var fees map[string]Fee

//var done chan struct{}

func main() {
	//log.SetFlags(0)
	initConfig()
	initSpotConnection()

	fees = make(map[string]Fee)
	symbols = make(map[string]*SymbolInfo)
	initFeesMap()
	initCurrencyGraph()

	asset := "BTC"
	log.Println("Using asset:", asset)
	log.Printf("Searching cycles with %d deals...", cycleDepth)
	findCycles(currencyGraph.nodes[asset])
	log.Printf("%d cycles has been found\n\n", len(cycles))

	done := make(chan struct{})

	log.Println("Listening wsBookTicker...")
	go wsBookTicker(done, asset)

	<-done
}
