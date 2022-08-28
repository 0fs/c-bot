package main

import (
	"context"
	"github.com/0fs/c-bot/utils"
	"github.com/adshao/go-binance/v2"
	"log"
)

var spotClient *binance.Client

var symbols map[string]*SymbolInfo

type SymbolInfo struct {
	symbol binance.Symbol
	book   *binance.WsBookTickerEvent
}

func initSpotConnection() {
	log.Println("Initializing binance spot connection...")
	binance.UseTestnet = config.GetBool("binance.test")
	spotClient = binance.NewClient(config.GetString("binance.api.spot.key"), config.GetString("binance.api.spot.secret"))

	timeOffset, err := spotClient.NewSetServerTimeService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	spotClient.TimeOffset = timeOffset

	updateBalances()

	log.Println("Done.")
	log.Println()
}

func updateBalances() {
	account, err := spotClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Non-zero balances:")
	for _, balance := range account.Balances {
		free := utils.Stf(balance.Free)
		locked := utils.Stf(balance.Locked)
		if free > 0 || locked > 0 {
			log.Printf("Free: %s %s Locked: %s %s ", balance.Free, balance.Asset, balance.Locked, balance.Asset)
		}
	}
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
