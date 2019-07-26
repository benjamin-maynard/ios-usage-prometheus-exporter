package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

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
