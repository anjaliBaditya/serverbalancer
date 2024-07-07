package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	config, err := LoadConfig("config.json")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	serverList := []Server{}
	for _, addr := range config.Servers {
		serverList = append(serverList, NewSimpleServer(addr))
	}

	lb := NewLoadBalancer(config.Port, serverList)
	go lb.HealthCheckServers(time.Duration(config.HealthCheckInterval) * time.Second)

	http.HandleFunc("/", lb.HandleRedirect)

	fmt.Printf("Load Balancer started at :%s\n", lb.port)
	err = http.ListenAndServe(":"+lb.port, nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
