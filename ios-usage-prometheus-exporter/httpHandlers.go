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

	log.Println("[Response: 404 Not Found] [Client IP: " + r.RemoteAddr + "] [Request URI: " + r.RequestURI + "]")

	response := responseMessage{"error", "Page not found"}
	js, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write(js)
	return

}

// HTTP Handler for totalAppOpens
func totalAppOpensHandler(w http.ResponseWriter, r *http.Request) {

	apiKey, apiKeyOk := r.URL.Query()["apiKey"]
	deviceName, deviceNameOk := r.URL.Query()["deviceName"]
	appName, appNameOk := r.URL.Query()["appName"]

	var response responseMessage
	var missingParam bool

	if !apiKeyOk || len(apiKey[0]) < 1 {

		log.Println("[Action: incTotalAppOpens] [Response: 400 Bad Request (Missing apiKey)] [Client IP: " + r.RemoteAddr + "]")
		response = responseMessage{"error", "The apiKey Query Parameter was missing from the request"}
		missingParam = true

	} else if !deviceNameOk || len(deviceName[0]) < 1 {

		log.Println("[Action: incTotalAppOpens] [Response: 400 Bad Request (Missing deviceName)] [Client IP: " + r.RemoteAddr + "]")
		response = responseMessage{"error", "The deviceName Query Parameter was missing from the request"}
		missingParam = true

	} else if !appNameOk || len(appName[0]) < 1 {

		log.Println("[Action: incTotalAppOpens] [Response: 400 Bad Request (Missing appName)] [Client IP: " + r.RemoteAddr + "]")
		response = responseMessage{"error", "The appName Query Parameter was missing from the request"}
		missingParam = true

	}

	if missingParam == true {

		js, _ := json.Marshal(response)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return

	}

	if apiKey[0] != registeredAPIKey {

		log.Println("[Action: incTotalAppOpens] [Response: 401 Unauthorized] [Client IP: " + r.RemoteAddr + "]")
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

	// Increment Counter and Return Response
	appOpens.With(prometheus.Labels{"appName": appName[0], "deviceName": deviceName[0]}).Inc()

	log.Println("[Action: incTotalAppOpens] [Response: 200 OK] [Client IP: " + r.RemoteAddr + "] [Device: " + string(deviceName[0]) + "] [App: " + string(appName[0]) + "]")
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
