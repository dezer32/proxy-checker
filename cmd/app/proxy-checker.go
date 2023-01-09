package main

import (
	"flag"
	"fmt"
	"github.com/dezer32/proxy-checker/pkg/checker"
	"time"
)

var (
	inputFileName, outputFileName string
)

func init() {
	defaultOutputFileName := fmt.Sprintf("proxies.checked.%d.json", time.Now().Unix())
	flag.StringVar(&inputFileName, "i", "proxies.json", "Path to file with proxies.")
	flag.StringVar(&outputFileName, "o", defaultOutputFileName, "Path to file with checked json.")
	flag.Parse()

	checker.LoadConfig()
}

func main() {
	checker.Run(inputFileName, outputFileName)
}
