package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type LoadBalancer struct {
	port            string
	roundRobinIndex int
	serverList      []Server
	mutex           sync.Mutex
}

func NewLoadBalancer(port string, serverList []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinIndex: 0,
		serverList:      serverList,
	}
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()

	for {
		server := lb.serverList[lb.roundRobinIndex%len(lb.serverList)]
		lb.roundRobinIndex++
		if server.IsAlive() {
			fmt.Printf("Selected server: %s\n", server.Address())
			return server
		}
	}
}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request: %s\n", req.URL.Path)
	targetServer := lb.getNextAvailableServer()
	targetServer.Serve(rw, req)
}

func (lb *LoadBalancer) HandleRedirect(rw http.ResponseWriter, req *http.Request) {
	lb.ServeProxy(rw, req)
}

func (lb *LoadBalancer) HealthCheckServers(interval time.Duration) {
	for {
		time.Sleep(interval)
		lb.mutex.Lock()
		for i, server := range lb.serverList {
			if !server.IsAlive() {
				fmt.Printf("Server %s is down. Removing from list.\n", server.Address())
				lb.serverList = append(lb.serverList[:i], lb.serverList[i+1:]...)
			}
		}
		lb.mutex.Unlock()
	}
}
