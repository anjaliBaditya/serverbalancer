package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type SimpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

type Server interface {
	Address() string
	IsAlive() bool
	Serve(http.ResponseWriter, *http.Request)
}

func NewSimpleServer(addr string) *SimpleServer {
	serverURL, err := url.Parse(addr)
	if err != nil {
		fmt.Printf("Failed to parse server URL: %v\n", err)
		return nil
	}

	return &SimpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverURL),
	}
}

func (s *SimpleServer) Address() string {
	return s.addr
}

func (s *SimpleServer) IsAlive() bool {
	resp, err := http.Get(s.addr)
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func (s *SimpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Forwarding request to %s\n", s.addr)
	s.proxy.ServeHTTP(rw, req)
}
