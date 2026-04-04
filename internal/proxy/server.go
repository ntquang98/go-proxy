// Package proxy
package proxy

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/elazarl/goproxy"
	"github.com/ntquang98/go-proxy/internal/config"
	"github.com/ntquang98/go-proxy/internal/rules"
)

func NewServer(cfg *config.Config, engine *rules.Engine) *http.Server {
	proxy := goproxy.NewProxyHttpServer()

	caCert, err := os.ReadFile(cfg.Cert.CertPath)
	if err != nil {
		log.Fatal(err)
	}

	caKey, err := os.ReadFile(cfg.Cert.CertKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	cert, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		log.Fatal(err)
	}

	goproxy.GoproxyCa = cert

	handler := NewHandler(engine)

	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(handler.HandleRequest)
	proxy.OnResponse().DoFunc(handler.HandleResponse)
	// proxy.Verbose = true

	return &http.Server{
		Addr:    cfg.Proxy.Addr,
		Handler: proxy,
	}
}
