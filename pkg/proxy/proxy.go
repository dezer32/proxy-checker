package proxy

import (
	"fmt"
	"github.com/gookit/ini/v2"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Proxy struct {
	Ip        string `json:"ip"`
	Port      uint   `json:"port"`
	Protocol  string `json:"protocol"`
	Country   string `json:"country"`
	IsWorking bool   `json:"is_working"`
}

func (p *Proxy) HealthCheck() {
	proxyUrl, err := url.Parse(fmt.Sprintf("%s://%s:%d", strings.ToLower(p.Protocol), p.Ip, p.Port))
	if err != nil {
		p.IsWorking = false
		return
	}

	timeout, err := time.ParseDuration(ini.String("timeout", ""))
	if err != nil {
		timeout = 60 * time.Second
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyUrl),
			IdleConnTimeout: timeout,
		},
	}

	log.Printf("Health check for '%s'.", proxyUrl)
	resp, err := client.Get(ini.String("check_url"))
	if err != nil {
		log.Printf("Proxy: '%s', error: '%s'.", proxyUrl, err)
		return
	}
	defer resp.Body.Close()

	p.IsWorking = resp.StatusCode == 200
	log.Printf("Proxy: %s, check: %v", proxyUrl, p.IsWorking)
}
