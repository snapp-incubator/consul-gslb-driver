package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func RunServer(cancel context.CancelFunc, listen string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	startServer(cancel, listen, mux)
}

func startServer(cancel context.CancelFunc, listen string, mux *http.ServeMux) {
	srv := &http.Server{
		Addr:    listen,
		Handler: mux,
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Listening on port %s ...\n", listen)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	<-sigCh
	cancel()

	timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelTimeout()

	if err := srv.Shutdown(timeoutCtx); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
}
