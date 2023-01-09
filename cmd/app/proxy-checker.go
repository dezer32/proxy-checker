package main

import (
	"flag"
	"fmt"
	"github.com/dezer32/proxy-checker/pkg/proxy"
	"github.com/gookit/ini/v2"
	"log"
	"sync"
	"time"
)

var (
	mutex          = sync.Mutex{}
	wg             = sync.WaitGroup{}
	pCh            = make(chan proxy.Proxy)
	checkedProxies = proxy.Proxies{}

	inputFileName, outputFileName string
)

func init() {
	defaultOutputFileName := fmt.Sprintf("proxies.checked.%d.json", time.Now().Unix())
	flag.StringVar(&inputFileName, "i", "proxies.json", "Path to file with proxies.")
	flag.StringVar(&outputFileName, "o", defaultOutputFileName, "Path to file with checked json.")
	flag.Parse()

	loadConfig()
}

func main() {
	proxies := proxy.Proxies{}
	proxies.Load(inputFileName)
	wg.Add(len(proxies.List))

	go runCheck(proxies.List)
	go consumeChecked()

	wg.Wait()

	if len(checkedProxies.List) > 0 {
		checkedProxies.Save(outputFileName)
	} else {
		log.Println("All proxies is dead =(")
	}
}

func runCheck(list []proxy.Proxy) {
	for _, p := range list {
		go func(proxy proxy.Proxy, pCh chan proxy.Proxy) {
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
		mutex.Lock()
		checkedProxies.List = append(checkedProxies.List, p)
		mutex.Unlock()
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
