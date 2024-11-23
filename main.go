package main

import (
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	addr := ":8080"
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("POST /v1/tester", testerHandler)

	log.Infof("dummy-tester is listening at %s...", addr)
	log.Error(http.ListenAndServe(addr, mux))
}

func testerHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ExpectedStatus  int    `json:"expectedStatus"`
		RequestDuration string `json:"requestDuration"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.WithFields(log.Fields{"err": err}).Error("Failed to decode request body")
		return
	}

	// Log the received values
	log.Infof("Received expectedStatus: %d, requestDuration: %s", body.ExpectedStatus, body.RequestDuration)

	// Parse the request duration
	duration, err := time.ParseDuration(body.RequestDuration)
	if err != nil {
		http.Error(w, "Invalid request duration", http.StatusBadRequest)
		log.WithFields(log.Fields{"err": err}).Error("Failed to parse request duration")
		return
	}

	// Wait for the specified duration
	time.Sleep(duration)

	// Return the expected status code
	w.WriteHeader(body.ExpectedStatus)
	w.Write([]byte("Request processed"))
	log.Infoln("Request processed")
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := make(map[string]string)
	resp["message"] = "Status OK"

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to marshal message")
	}
	return
}
