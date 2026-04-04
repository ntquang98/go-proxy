package main

import (
	"fmt"
	"log"

	"github.com/ntquang98/go-proxy/internal/config"
	"github.com/ntquang98/go-proxy/internal/proxy"
	"github.com/ntquang98/go-proxy/internal/rules"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	rulesList, err := rules.LoadRules("configs/rules.json")
	if err != nil {
		log.Fatal(err)
	}

	engine := rules.NewEngine(rulesList)
	server := proxy.NewServer(cfg, engine)

	fmt.Println("Proxy running on ", cfg.Proxy.Addr)
	log.Fatal(server.ListenAndServe())
}
