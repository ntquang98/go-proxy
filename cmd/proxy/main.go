package main

import (
	"log"

	"github.com/ntquang98/go-proxy/internal/proxy"
	"github.com/ntquang98/go-proxy/internal/rules"
)

func main() {
	rulesList, err := rules.LoadRules("configs/rules.json")
	if err != nil {
		log.Fatal(err)
	}

	engine := rules.NewEngine(rulesList)
	server := proxy.NewServer(engine)

	log.Println("Proxy running on :3333")
	log.Fatal(server.ListenAndServe())
}
