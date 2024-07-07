package main

import (
	"fmt"
	"net/http"
)

func main() {
	serverList := []Server{
		NewSimpleServer("https://www.facebook.com/"),
		NewSimpleServer("https://www.bing.com/"),
		NewSimpleServer("https://www.duckduckgo.com/"),
	}

	lb := NewLoadBalancer("8009", serverList)
	http.HandleFunc("/", lb.HandleRedirect)

	fmt.Printf("Load Balancer started at :%s\n", lb.port)
	err := http.ListenAndServe(":"+lb.port, nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
