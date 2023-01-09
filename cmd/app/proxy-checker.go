package main

import (
	"flag"
	"github.com/dezer32/proxy-checker/pkg"
	"github.com/gookit/ini/v2"
	"log"
	"sync"
)

var (
	wg             = sync.WaitGroup{}
	pCh            = make(chan pkg.Proxy)
	checkedProxies = pkg.Proxies{}

	inputFileName, outputFileName string
)

func init() {
	flag.StringVar(&inputFileName, "i", "proxies.json", "Path to file with proxies.")
	flag.StringVar(&outputFileName, "o", "proxies.checked.json", "Path to file with checked json.")
	flag.Parse()

	loadConfig()
}

func main() {
	proxies := pkg.Proxies{}
	proxies.Load(inputFileName)
	wg.Add(len(proxies.List))

	go runCheck(proxies.List)
	go consumeChecked()

	wg.Wait()

	checkedProxies.Save(outputFileName)
}

func runCheck(list []pkg.Proxy) {
	for _, p := range list {
		go func(proxy pkg.Proxy, pCh chan pkg.Proxy) {
			defer wg.Done()
			proxy.HealthCheck()
			if proxy.IsWorking {
				wg.Add(1)
				pCh <- proxy
			}
		}(p, pCh)
	}
}

func consumeChecked() {
	//d, err := time.ParseDuration(config.String("timeout", ""))
	//if err != nil {
	//	d = 60 * time.Second
	//}
	//
	//timeoutCh := time.After(d * 2)

	for p := range pCh {
		checkedProxies.List = append(checkedProxies.List, p)
		wg.Done()
	}
}

func loadConfig() {
	ini.WithOptions(ini.ParseEnv)

	err := ini.LoadFiles("configs/check.ini")
	if err != nil {
		log.Fatalf("Config don't load. Error: '%s'.", err)
	}
}
