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
	fee      float64
}

type Currency struct {
	value string
	books *binance.WsBookTickerEvent
}

type Fee struct {
	m float64 // Maker
	t float64 // Taker
}

var color map[*Currency]struct{}
var cycles map[int][]*Currency
var p map[*Currency]*Currency
var ncycle int

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
		fee:      fees[symbol].m,
	})
	g.edges[*n2] = append(g.edges[*n2], &Edge{
		symbol:   symbol,
		currency: n1,
		fee:      fees[symbol].t,
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
		s += node.String()
		near := g.edges[*node]
		for j, item := range near {
			if j != 0 {
				for i := 0; i < len(node.String()); i++ {
					s += " "
				}
			}
			s += fmt.Sprintf(" -> %s - Symbol: %s, Fee: %.1f%%", item.currency.String(), item.symbol, item.fee*100)
			fmt.Println(s)
			s = ""
		}
	}
	g.lock.RUnlock()
}

func initFeesMap() {
	log.Println("Initializing fees map...")

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

		if binance.SymbolStatusType(symbol.Status) == binance.SymbolStatusTypeTrading {
			currencyGraph.AddEdge(symbol.Symbol, existNodes[symbol.BaseAsset], existNodes[symbol.QuoteAsset])

			if _, ok := symbols[symbol.Symbol]; !ok {
				symbols[symbol.Symbol] = &SymbolInfo{
					symbol: symbol,
					book:   nil,
				}
			}
		}
	}

	log.Println("Done.")
	log.Println()
}

func findCycles(currency *Currency) {
	ncycle = 0
	color = make(map[*Currency]struct{})
	cycles = make(map[int][]*Currency)
	p = make(map[*Currency]*Currency)
	dfs(currency, 0)
	delete(cycles, ncycle) // TODO: check
}

func saveCycle(f *Currency, s *Currency) int {
	delete(cycles, ncycle)
	cycles[ncycle] = append(cycles[ncycle], s)
	for v := f; v != s; v = p[v] {
		cycles[ncycle] = append(cycles[ncycle], v)
	}
	cycles[ncycle] = append(cycles[ncycle], s)

	return len(cycles[ncycle])
}

func dfs(c *Currency, depth int) {
	depth++
	if depth > cycleDepth {
		return
	}
	color[c] = struct{}{}
	for _, e := range currencyGraph.edges[*c] {
		to := e.currency
		if _, ok := color[to]; !ok {
			p[to] = c
			dfs(to, depth)
		} else {
			if saveCycle(c, to) > cycleDepth {
				ncycle++
			}
		}
	}
	delete(color, c)
}
