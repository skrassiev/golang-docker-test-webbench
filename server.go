package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var (
	abort bool
)

const (
	SOCK = "/tmp/go.sock"
)

type Server struct {
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.RequestURI != "/silent" {
		if r.RequestURI == "/exit" {
			os.Exit(2)
		}
		log.Printf("Request from [%v] -> %v\n", r.RemoteAddr, r.RequestURI)
	}

	body := fmt.Sprintf("Hello World at %v\n", time.Now())
	// Try to keep the same amount of headers
	w.Header().Set("Server", "gophr")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", fmt.Sprint(len(body)))
	fmt.Fprint(w, body)
}

func main() {

	runtime.GOMAXPROCS(8)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	signal.Notify(sigchan, syscall.SIGTERM)

	server := Server{}

	go func() {
		http.Handle("/", server)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		tcp, err := net.Listen("tcp", ":9001")
		if err != nil {
			log.Fatal(err)
		}
		fcgi.Serve(tcp, server)
	}()

	go func() {
		unix, err := net.Listen("unix", SOCK)
		if err != nil {
			log.Fatal(err)
		}
		fcgi.Serve(unix, server)
	}()

	fmt.Println("HTTP server is running on port 8080. Try 'curl -v localhost:8080' from another windows. To terminate, call 'curl localhost:8080/exit'")

	<-sigchan

	if err := os.Remove(SOCK); err != nil {
		log.Fatal(err)
	}
}
