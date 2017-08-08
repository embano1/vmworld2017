package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	addr    = ":8080"
	healthZ = "/healthz"
	preStop = "/prestop"
)

var (
	version string
	build   string
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc(healthZ, healthz)
	mux.HandleFunc(preStop, prestop)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		h, _ := os.Hostname()
		log.Printf("Agent: %s, Source: %s, Path: %s\n", r.UserAgent(), r.RemoteAddr, r.RequestURI)
		fmt.Fprintf(w, "Hello Gopher!\n")
		fmt.Fprintf(w, "Hostname: %s\n", h)
		fmt.Fprintf(w, "Version: %s\n", version)
		fmt.Fprintf(w, "Build: %s\n", build)
	})

	// register pprof handlers
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	srv := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// set up root context (http.Shutdown()) and prepare to catch OS signals
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	defer func() {
		signal.Stop(c)
	}()

	// launch httpd in a separate goroutine
	log.Println("Starting webserver")
	log.Printf("Additional handlers: %v, %v, /debug/pprof", healthZ, preStop)
	go srv.ListenAndServe()

	// catch OS sigs
	sig := <-c
	log.Printf("Got %v\n", sig)
	log.Println("Attempting graceful shutdown (closing open handlers, etc.)")
	err := srv.Shutdown(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	} else {
		log.Println("Done")
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func prestop(w http.ResponseWriter, r *http.Request) {
	// Gives Kubernetes time to remove this container from endpoint (service) list before sending SIGTERM
	// More details here: https://github.com/kubernetes/ingress/issues/322#issuecomment-315884499
	time.Sleep(10 * time.Second)
}
