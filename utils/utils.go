package utils

import (
	"fmt"
	"log"
	"strconv"
)

func Stf(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
	}

	return f
}

func Fts(f float64) string {
	return fmt.Sprintf("%.8f", f)
}

func IntMax(a, b int) int {
	if a > b {
		return a
	}
	return b
}
