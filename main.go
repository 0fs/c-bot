package main

import (
	"github.com/adshao/go-binance/v2"
	"log"
)

var currencyGraph CurrencyGraph
var fees map[string]Fee
var symbols map[string]*SymbolInfo

type SymbolInfo struct {
	symbol binance.Symbol
	book   *binance.WsBookTickerEvent
}

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

func wsBookTicker(done chan struct{}) {
	wsHandler := func(event *binance.WsBookTickerEvent) {
		if _, ok := symbols[event.Symbol]; ok {
			symbols[event.Symbol].book = event
		}
	}

	errHandler := func(err error) {
		log.Fatal(err)
	}

	doneC, _, err := binance.WsAllBookTickerServe(wsHandler, errHandler)
	if err != nil {
		log.Fatal(err)
		return
	}

	msg := <-doneC
	done <- msg
}
