package main

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"log"
	"sync"
)

type CurrencyGraph struct {
	nodes []*Currency
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

// AddNode adds a node to the graph
func (g *CurrencyGraph) AddNode(n *Currency) {
	g.lock.Lock()
	g.nodes = append(g.nodes, n)
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
	for i := 0; i < len(g.nodes); i++ {
		s += g.nodes[i].String() + " -> "
		near := g.edges[*g.nodes[i]]
		for j := 0; j < len(near); j++ {
			s += near[j].currency.String() + "(" + near[j].symbol + ") "
		}
		s += "\n"
	}
	fmt.Println(s)
	g.lock.RUnlock()
}

func initCurrencyGraph() {
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
}