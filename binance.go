package main

import (
	"context"
	"github.com/0fs/c-bot/utils"
	"github.com/adshao/go-binance/v2"
	"log"
)

var spotClient *binance.Client

var btcBalance binance.Balance
var usdtBalance binance.Balance

var symbol string
var limit int
var interval string
var qty string

var firstTrade = true

func initSpotConnection() {
	log.Println("Spot initialization...")
	binance.UseTestnet = config.GetBool("binance.test")
	spotClient = binance.NewClient(config.GetString("binance.api.spot.key"), config.GetString("binance.api.spot.secret"))

	timeOffset, err := spotClient.NewSetServerTimeService().Do(context.Background())
	if err != nil {
		log.Fatal("Could not get server time offset")
	}
	spotClient.TimeOffset = timeOffset

	symbol = config.GetString("binance.symbol")
	limit = config.GetInt("binance.limit")
	interval = config.GetString("binance.interval")
	qty = config.GetString("binance.qty")

	updateBalances()
}

func updateBalances() {
	account, err := spotClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, balance := range account.Balances {
		switch balance.Asset {
		case "BTC":
			btcBalance = balance
			break
		case "USDT":
			usdtBalance = balance
			break
		default:
			break
		}
	}

	log.Printf("BTC: %.8f | USDT %.8f", utils.Stf(btcBalance.Free), utils.Stf(usdtBalance.Free))
}
