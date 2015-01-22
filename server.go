package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

var (
	abort bool
)

const (
	VERSION = 127
)

type Server struct {
	ipAddress string
}

func (s *Server) ip() string {
	if len(s.ipAddress) == 0 {
		s.ipAddress = "<unknown>"
		addrs, err := net.InterfaceAddrs()

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, address := range addrs {

			// check the address type and if it is not a loopback the display it
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					s.ipAddress = ipnet.IP.String()
				}

			}
		}
	}
	return s.ipAddress
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.RequestURI != "/silent" {
		if r.RequestURI == "/exit" {
			os.Exit(2)
		}

		log.Printf("Request from [%v] -> %v\n", r.RemoteAddr, r.RequestURI)

		if v, err := strconv.Atoi(r.RequestURI[1:]); err == nil {
			log.Printf("Sleeping for %v seconds\n", v)
			time.Sleep(time.Duration(v) * time.Second)
		}
	}

	body := fmt.Sprintf("[%v v:%v] Hello World at %v\n", s.ip(), VERSION, time.Now())
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

		srv := http.Server{
			Addr:        ":8080",
			Handler:     server,
			ReadTimeout: 2 * time.Second,
		}

		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("HTTP server is running on port 8080. Try 'curl -v localhost:8080' from another windows. To terminate, call 'curl localhost:8080/exit'")

	s := <-sigchan
	fmt.Println("Handling received signal ", s)
}
