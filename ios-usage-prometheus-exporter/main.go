package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

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

	// Check and set PROMETHEUS_PORT Environment Variable
	promPort := os.Getenv("PROMETHEUS_PORT")
	if len(promPort) == 0 {
		log.Fatal("The PROMETHEUS_PORT Environment Variable was not set")
	} else if _, err := strconv.Atoi(promPort); err != nil {
		log.Fatal("The PROMETHEUS_PORT Environment Variable should only contain numbers")
	}

	// Check and set WEBSERVER_PORT Environment Variable
	wsPort := os.Getenv("WEBSERVER_PORT")
	if len(wsPort) == 0 {
		log.Fatal("The WEBSERVER_PORT Environment Variable was not set")
	} else if _, err := strconv.Atoi(wsPort); err != nil {
		log.Fatal("The WEBSERVER_PORT Environment Variable should only contain numbers")
	}

	// Setup HTTP Server for Prometheus Metrics Endpoint. Run it as a Goroutine so it doesn't block
	promServer := http.NewServeMux()
	promServer.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":"+promPort, promServer)

	// Setup HTTP Server for API to Receive Metrics
	reportingServer := http.NewServeMux()
	reportingServer.HandleFunc("/", defaultHandler)
	reportingServer.HandleFunc("/api/v1.0/incTotalAppOpens/", totalAppOpensHandler)
	http.ListenAndServe(":"+wsPort, reportingServer)

}
