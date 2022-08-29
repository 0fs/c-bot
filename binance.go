package main

import (
	"context"
	"errors"
	"fmt"
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

func wsBookTicker(done chan struct{}, asset string) {
	wsHandler := func(event *binance.WsBookTickerEvent) {
		if _, ok := symbols[event.Symbol]; ok {
			symbols[event.Symbol].book = event

			if symbols[event.Symbol].symbol.BaseAsset != asset && symbols[event.Symbol].symbol.QuoteAsset != asset {
				for _, cycle := range cycles {
					profit, prices, err := checkCycle(cycle)
					if err != nil {
						continue
					}
					profit = (profit - 1.0) * 100.0
					if profit > 0 {
						rstr := ""
						for i, c := range cycle {
							if i == cycleDepth {
								rstr += fmt.Sprintf("%s expected profit:  %.5f%%", c.value, profit)
							} else {
								rstr += fmt.Sprintf("%s - %.8f - ", c.value, prices[i])
							}
						}
						log.Println(rstr)
						//log.Printf("Cycle: %v expected profit: %.5f%%", cycle, profit)
					}
				}
			}
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

func checkCycle(steps []*Currency) (float64, []float64, error) {
	var price float64
	prices := make([]float64, cycleDepth)
	base, qty := 10000.0, 10000.0
	for i := 1; i < len(steps); i++ {
		edges := currencyGraph.edges[*steps[i-1]]
		for _, edge := range edges {
			if edge.currency == steps[i] {
				if symbols[edge.symbol].book == nil {
					return 0, nil, errors.New("no book")
				}

				if edge.currency.value == symbols[edge.symbol].symbol.BaseAsset {
					// Если следующий ассет базовый - берем минимальную цену продажи BTC за USDT
					// USDT -> BTC (BTCUSDT) BTC -> базовый
					price = utils.Stf(symbols[edge.symbol].book.BestAskPrice)
				} else {
					// Максимальная цена покупки BTC за USDT
					// BTC -> USDT (BTCUSDT) BTC -> базовый
					price = utils.Stf(symbols[edge.symbol].book.BestBidPrice)
				}

				prices[i-1] = price

				if edge.currency.value == symbols[edge.symbol].symbol.BaseAsset {
					qty /= price
				} else {
					qty *= price
				}
				qty -= qty * edge.fee
			}
		}
	}

	return qty / base, prices, nil
}

func wsDepth(done chan struct{}, symbol string) {
	wsHandler := func(event *binance.WsDepthEvent) {
		log.Println("ASKS: ", event.Asks[:3])
		log.Println("BIDS: ", event.Bids[:3])
	}

	errHandler := func(err error) {
		log.Fatal(err)
	}

	doneC, _, err := binance.WsDepthServe(symbol, wsHandler, errHandler)
	if err != nil {
		log.Fatal(err)
		return
	}

	msg := <-doneC
	done <- msg
}
