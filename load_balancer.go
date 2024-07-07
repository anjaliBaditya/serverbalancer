package main

import (
	"fmt"
	"net/http"
)

type LoadBalancer struct {
	port            string
	roundRobinIndex int
	serverList      []Server
}

// NewLoadBalancer creates a new instance of LoadBalancer
func NewLoadBalancer(port string, serverList []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinIndex: 0,
		serverList:      serverList,
	}
}

// getNextAvailableServer returns the next available server in a round-robin manner
func (lb *LoadBalancer) getNextAvailableServer() Server {
	for {
		server := lb.serverList[lb.roundRobinIndex%len(lb.serverList)]
		lb.roundRobinIndex++
		if server.IsAlive() {
			fmt.Printf("Selected server: %s\n", server.Address())
			return server
		}
	}
}

// ServeProxy handles the request and forwards it to the target server
func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request: %s\n", req.URL.Path)
	targetServer := lb.getNextAvailableServer()
	targetServer.Serve(rw, req)
}

// HandleRedirect is the HTTP handler for redirecting requests to the target server
func (lb *LoadBalancer) HandleRedirect(rw http.ResponseWriter, req *http.Request) {
	lb.ServeProxy(rw, req)
}
