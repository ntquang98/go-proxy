package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log/slog"

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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	engine := rules.NewEngine(rulesList)
	server := proxy.NewServer(cfg, engine)

	go func() {
		log.Println("[START] proxy running on", cfg.Proxy.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	log.Println("received signal, shutting down...", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[SHUTDONW] timed out, forced shutdown: %v\n", err)
	} else {
		log.Println("[SHUTDOWN] clean")
	}
}
