package main

import (
	"github.com/adshao/go-binance/v2"
	"log"
)

var books map[string]*binance.WsBookTickerEvent

var currencyGraph CurrencyGraph
var fees map[string]Fee

//var done chan struct{}

func main() {
	log.SetFlags(0)
	fees = make(map[string]Fee)
	initConfig()
	initSpotConnection()
	initFeesMap()
	initCurrencyGraph()
	currencyGraph.String()
	return

	done := make(chan struct{})

	go wsBookTicker(done, "BTCUSDT", "ETHUSDT")

	<-done
}

func wsBookTicker(done chan struct{}, symbols ...string) {
	books := make(map[string]*binance.WsBookTickerEvent)
	for _, symbol := range symbols {
		books[symbol] = nil
	}

	wsHandler := func(event *binance.WsBookTickerEvent) {
		if _, ok := books[event.Symbol]; ok {
			books[event.Symbol] = event
			log.Printf("%s : Buy: %s x %s - Sell: %s x %s\n", event.Symbol, event.BestAskPrice, event.BestAskQty, event.BestBidPrice, event.BestBidQty)
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
