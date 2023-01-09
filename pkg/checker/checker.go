package checker

import (
	"github.com/dezer32/proxy-checker/pkg/proxy"
	"github.com/gookit/ini/v2"
	"log"
	"sync"
)

var (
	mutex          = sync.Mutex{}
	wg             = sync.WaitGroup{}
	checkedProxies = proxy.Proxies{}
	pCh            = make(chan proxy.Proxy)
)

func Run(inputFileName string, outputFileName string) {
	proxies := proxy.Proxies{}
	proxies.Load(inputFileName)
	wg.Add(len(proxies.List))

	go runCheck(proxies.List)
	go consumeChecked()

	wg.Wait()

	checkedProxies.Save(outputFileName)
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
	for p := range pCh {
		mutex.Lock()
		checkedProxies.List = append(checkedProxies.List, p)
		mutex.Unlock()
		wg.Done()
	}
}

func LoadConfig() {
	ini.WithOptions(ini.ParseEnv)

	err := ini.LoadFiles("configs/check.ini")
	if err != nil {
		log.Fatalf("Config don't load. Error: '%s'.", err)
	}
}
