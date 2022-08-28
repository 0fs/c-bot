package main

import (
	"context"
	"fmt"
	"github.com/0fs/c-bot/utils"
	"github.com/adshao/go-binance/v2"
	"log"
	"sync"
)

type CurrencyGraph struct {
	nodes map[string]*Currency
	edges map[Currency][]*Edge
	lock  sync.RWMutex
}

type Edge struct {
	symbol   string
	currency *Currency
}

type Currency struct {
	value string
	books *binance.WsBookTickerEvent
}

var fees map[string]Fee

type Fee struct {
	m float64 // Maker
	t float64 // Taker
}

// AddNode adds a node to the graph
func (g *CurrencyGraph) AddNode(n *Currency) {
	g.lock.Lock()
	if g.nodes == nil {
		g.nodes = make(map[string]*Currency)
	}

	if _, ok := g.nodes[n.value]; !ok {
		g.nodes[n.value] = n
	}
	g.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *CurrencyGraph) AddEdge(symbol string, n1, n2 *Currency) {
	g.lock.Lock()
	if g.edges == nil {
		g.edges = make(map[Currency][]*Edge)
	}
	g.edges[*n1] = append(g.edges[*n1], &Edge{
		symbol:   symbol,
		currency: n2,
	})
	g.edges[*n2] = append(g.edges[*n2], &Edge{
		symbol:   symbol,
		currency: n1,
	})
	g.lock.Unlock()
}

func (n *Currency) String() string {
	return n.value
}

func (g *CurrencyGraph) String() {
	g.lock.RLock()
	s := ""
	for _, node := range g.nodes {
		s += node.String() + " -> "
		near := g.edges[*node]
		for _, item := range near {
			s += item.currency.String() + "(" + item.symbol + ") "
		}
		s += "\n"
	}
	fmt.Println(s)
	g.lock.RUnlock()
}

func initFeesMap() {
	log.Println("Initializing fees map...")
	fees := make(map[string]Fee)

	rs, err := spotClient.NewTradeFeeService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range rs {
		fees[item.Symbol] = Fee{
			m: utils.Stf(item.MakerCommission),
			t: utils.Stf(item.TakerCommission),
		}
	}

	log.Println("Done.")
	log.Println()
}

func initCurrencyGraph() {
	log.Println("Initializing currency graph...")
	rs, err := spotClient.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	existNodes := make(map[string]*Currency)

	for _, symbol := range rs.Symbols {

		if _, ok := existNodes[symbol.BaseAsset]; !ok {
			existNodes[symbol.BaseAsset] = &Currency{value: symbol.BaseAsset}
			currencyGraph.AddNode(existNodes[symbol.BaseAsset])
		}

		if _, ok := existNodes[symbol.QuoteAsset]; !ok {
			existNodes[symbol.QuoteAsset] = &Currency{value: symbol.QuoteAsset}
			currencyGraph.AddNode(existNodes[symbol.QuoteAsset])
		}

		currencyGraph.AddEdge(symbol.Symbol, existNodes[symbol.BaseAsset], existNodes[symbol.QuoteAsset])
	}

	log.Println("Done.")
	log.Println()
}
