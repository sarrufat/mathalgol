package main

import (
	"flag"
	"fmt"
	"mathalgol/Divisors/divisors"
	"os"
	"strings"
)

var index = 1

func handler(max int64, numd int) {
	fmt.Printf("%3d -> %d, %d div.\n", index, max, numd)
	index++
}

func parHandler(idx int, max int64, numd int) {
	fmt.Printf("%3d -> %d, %d div.\n", idx, max, numd)
}

func printCmd() {

	fmt.Println("cmd:", strings.Join(os.Args, " "))
}

func main() {
	nValue := flag.Int64("n", 100000, "top N")
	parBool := flag.Bool("par", false, "parallel")
	flag.Parse()
	printCmd()
	handler(1, 1)
	if *parBool {
		divisors.ParallelHighlyComposite(*nValue, parHandler)
		fmt.Printf("#discards %d, #discards2 %d, #discards3 %d\n", divisors.AtomicDiscards, divisors.AtomicDiscards2, divisors.AtomicDiscards3)
	} else {
		divisors.HighlyComposite(*nValue, handler)
	}
}
