package main

import "github.com/prometheus/client_golang/prometheus"

// Define Prometheus Metrics
var (
	appOpens = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ios_app_open_total",
			Help: "A metric to count how many times an application is opened on iOS, by device, by app.",
		},
		[]string{"appName", "deviceName"},
	)
)

// Register Prometheus Metrics
func init() {

	prometheus.MustRegister(appOpens)

}
