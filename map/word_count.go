package main

import (
	"golang.org/x/tour/wc"
	"strings"
)

func WordCount(s string) map[string]int {
	fields := strings.Fields(s)
	res := make(map[string]int)
	for _, e := range fields {
		elem, _ := res[e]
		res[e] = elem + 1
	}
	return res
}

func main() {
	wc.Test(WordCount)
}
