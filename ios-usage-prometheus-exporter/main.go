package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var registeredAPIKey string

func main() {

	// Check and set API_KEY Environment Variable
	registeredAPIKey = os.Getenv("API_KEY")
	if len(registeredAPIKey) == 0 {
		log.Fatal("The API_KEY Environment Variable was not set")
	} else if len(registeredAPIKey) < 10 {
		log.Fatal("The API_KEY is not a sufficient length, please ensure it is at least 10 characters long")
	}

	// Setup HTTP Server for Prometheus Metrics Endpoint. Run it as a Goroutine so it doesn't block
	promServer := http.NewServeMux()
	promServer.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":9001", promServer)

	// Setup HTTP Server for API to Receive Metrics
	reportingServer := http.NewServeMux()
	reportingServer.HandleFunc("/", defaultHandler)
	reportingServer.HandleFunc("/api/v1.0/totalAppOpens/", totalAppOpensHandler)
	http.ListenAndServe(":8080", reportingServer)

}
