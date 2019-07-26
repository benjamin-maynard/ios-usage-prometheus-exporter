package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

// Define responseMessage for JSON responses
type responseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Define  a default handler for 404 errors
func defaultHandler(w http.ResponseWriter, r *http.Request) {

	response := responseMessage{"error", "Page not found"}
	js, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(js)
	return

}

// HTTP Handler for totalAppOpens
func totalAppOpensHandler(w http.ResponseWriter, r *http.Request) {

	apiKey, ok := r.URL.Query()["apiKey"]
	deviceName, ok := r.URL.Query()["deviceName"]
	appName, ok := r.URL.Query()["appName"]

	var response responseMessage
	var missingParam bool

	if !ok || len(apiKey[0]) < 1 {

		log.Println("Error: apiKey Query Parameter missing from request")
		response = responseMessage{"error", "The apiKey Query Parameter was missing from the request"}
		missingParam = true

	} else if !ok || len(deviceName[0]) < 1 {

		log.Println("Error: deviceName Query Parameter missing from request")
		response = responseMessage{"error", "The deviceName Query Parameter was missing from the request"}
		missingParam = true

	} else if !ok || len(appName[0]) < 1 {

		log.Println("Error: appName Query Parameter missing from request")
		response = responseMessage{"error", "The appName Query Parameter was missing from the request"}
		missingParam = true

	}

	if missingParam == true {

		js, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return

	}

	if apiKey[0] != registeredAPIKey {

		log.Println("Error: Invalid API Key Specified")

		response := responseMessage{"error", "The apiKey specified was invalid"}

		js, err := json.Marshal(response)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(js)
		return

	}

	// Successful  Counter Increment
	log.Println("Incremented Prometheus Counter for " + string(appName[0]) + " on device: " + string(deviceName[0]) + ".")

	appOpens.With(prometheus.Labels{"appName": appName[0], "deviceName": deviceName[0]}).Inc()

	response = responseMessage{"success", "Incremented Prometheus Counter for " + string(appName[0]) + " on device: " + string(deviceName[0]) + "."}

	js, err := json.Marshal(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
	return

}
