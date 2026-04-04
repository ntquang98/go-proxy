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

// func main() {
// 	proxy := goproxy.NewProxyHttpServer()
//
// 	caCert, err := os.ReadFile("C:/Users/ADMIN/AppData/Local/mkcert/rootCA.pem")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	caKey, err := os.ReadFile("C:/Users/ADMIN/AppData/Local/mkcert/rootCA-key.pem")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	cert, err := tls.X509KeyPair(caCert, caKey)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	goproxy.GoproxyCa = cert
//
// 	proxy.Verbose = true
// 	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
//
// 	log.Println("Proxy running on :8080")
// 	log.Fatal(http.ListenAndServe("0.0.0.0:8080", proxy))
// }
