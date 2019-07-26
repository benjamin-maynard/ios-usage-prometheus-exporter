package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

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

	apiKey := r.Header.Get("apiKey")
	deviceName := r.Header.Get("deviceName")
	appName := r.Header.Get("appName")

	var response responseMessage
	var missingParam bool

	if len(apiKey) < 1 {

		log.Println("[Action: incTotalAppOpens] [Response: 400 Bad Request (Missing apiKey)] [Client IP: " + r.RemoteAddr + "]")
		response = responseMessage{"error", "The apiKey Header was missing from the request"}
		missingParam = true

	} else if len(deviceName) < 1 {

		log.Println("[Action: incTotalAppOpens] [Response: 400 Bad Request (Missing deviceName)] [Client IP: " + r.RemoteAddr + "]")
		response = responseMessage{"error", "The deviceName Header was missing from the request"}
		missingParam = true

	} else if len(appName) < 1 {

		log.Println("[Action: incTotalAppOpens] [Response: 400 Bad Request (Missing appName)] [Client IP: " + r.RemoteAddr + "]")
		response = responseMessage{"error", "The appName Header was missing from the request"}
		missingParam = true

	}

	if missingParam == true {

		js, _ := json.Marshal(response)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(js)
		return

	}

	if apiKey != registeredAPIKey {

		log.Println("[Action: incTotalAppOpens] [Response: 401 Unauthorized] [Client IP: " + r.RemoteAddr + "]")
		response := responseMessage{"error", "The apiKey specified was invalid"}

		js, _ := json.Marshal(response)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(js)
		return

	}

	// Ensure appName and deviceName are formatted nicely : Credit: https://golangcode.com/how-to-remove-all-non-alphanumerical-characters-from-a-string/
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Println("[Action: incTotalAppOpens] [Response: 500 Internal Server Error] [Client IP: " + r.RemoteAddr + "] [Error: Errror Removing Special Characters")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	appName = reg.ReplaceAllString(appName, "")
	deviceName = reg.ReplaceAllString(deviceName, "")

	// Increment Counter and Return Response
	appOpens.With(prometheus.Labels{"appName": appName, "deviceName": deviceName}).Inc()

	log.Println("[Action: incTotalAppOpens] [Response: 200 OK] [Client IP: " + r.RemoteAddr + "] [Device: " + string(deviceName) + "] [App: " + string(appName) + "]")
	response = responseMessage{"success", "Incremented Prometheus Counter for " + string(appName) + " on device: " + string(deviceName) + "."}

	js, err := json.Marshal(response)

	if err != nil {
		log.Println("[Action: incTotalAppOpens] [Response: 500 Internal Server Error] [Client IP: " + r.RemoteAddr + "] [Error: Errror Marshalling JSON Response")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)
	return

}
