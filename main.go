package main

import (
	"log"
)

var currencyGraph CurrencyGraph
var fees map[string]Fee

//var done chan struct{}

func main() {
	log.SetFlags(0)
	initConfig()
	initSpotConnection()

	fees = make(map[string]Fee)
	symbols = make(map[string]*SymbolInfo)

	initFeesMap()
	initCurrencyGraph()

	done := make(chan struct{})

	go wsBookTicker(done)

	<-done
}
