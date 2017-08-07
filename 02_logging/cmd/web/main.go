package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	version = "1.0"
	addr    = ":8080"
)

var (
	myLogger *log.Logger
)

func main() {

	logfile := flag.String("l", "/log/http.log", "Full path and name for logfile")
	flag.Parse()

	file, err := os.Create(*logfile)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	myLogger = log.New(file, "[http] ", log.LstdFlags|log.Lshortfile)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		h, _ := os.Hostname()
		myLogger.Printf("Agent: %s, Source: %s, Path: %s\n", r.UserAgent(), r.RemoteAddr, r.RequestURI)
		fmt.Fprintf(w, "Hello Gopher! I´ve been hardcoded to log to %v :(\n", *logfile)
		fmt.Fprintf(w, "Hostname: %s\n", h)
		fmt.Fprintf(w, "Version: %s\n", version)
	})

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// set up root context (http.Shutdown()) and prepare to catch OS signals
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	defer func() {
		// cancel()
		signal.Stop(c)
	}()

	myLogger.Println("Starting webserver")
	go srv.ListenAndServe()

	// catch OS sigs
	sig := <-c
	myLogger.Printf("Got %v\n", sig)
	myLogger.Println("Attempting graceful shutdown (closing open handlers, etc.)")
	err = srv.Shutdown(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	} else {
		myLogger.Println("Done")
	}
}
