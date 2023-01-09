package pkg

import (
	"encoding/json"
	"log"
	"os"
)

type Proxies struct {
	List []Proxy
}

func (p *Proxies) Load(fileName string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Proxy list file '%s' don't load. Error: %s", fileName, err)
	}

	var res []Proxy
	err = json.Unmarshal(data, &res)
	if err != nil {
		log.Fatalf("Json unmarshal error. Error: '%s'.", err)
	}

	p.List = append(p.List, res...)
}

func (p *Proxies) Save(fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Create file failed. Error: '%s'.", err)
	}
	defer file.Close()

	j, err := json.Marshal(p.List)
	if err != nil {
		log.Fatalf("Json marshal error. Error: '%s'.", err)
	}

	log.Printf("Save %d proxies to '%s'.", len(p.List), fileName)

	if _, err = file.Write(j); err != nil {
		log.Fatalf("Save proxies failed. Error: '%s'.", err)
	}
}

func (p *Proxies) fileExists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		return false
	}

	return true
}
