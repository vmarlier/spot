package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	// Initialize logrus
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// Define a flag for the port
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	addr := fmt.Sprintf(":%s", *port)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", logRequest(healthHandler))
	mux.HandleFunc("POST /v1/tester", logRequest(testerHandler))
	mux.HandleFunc("POST /v1/caller", logRequest(callerHandler))

	log.Infof("dummy-tester is listening at %s...", addr)
	log.Error(http.ListenAndServe(addr, mux))
}

func logRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.WithFields(log.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
		}).Info("Received request")

		handler(w, r)

		log.WithFields(log.Fields{
			"method":      r.Method,
			"url":         r.URL.String(),
			"elapsedTime": time.Since(start),
		}).Info("Handled request")
	}
}

func testerHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Handling /v1/tester request")

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
	log.Infof("Sleeping for %s", duration)
	time.Sleep(duration)

	// Return the expected status code
	w.WriteHeader(body.ExpectedStatus)
	w.Write([]byte("Request processed"))
	log.Infoln("Request processed")
}

func callerHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Handling /v1/caller request")

	var body struct {
		URL             string `json:"url"`
		Port            string `json:"port"`
		ExpectedStatus  int    `json:"expectedStatus"`
		RequestDuration string `json:"requestDuration"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.WithFields(log.Fields{"err": err}).Error("Failed to decode request body")
		return
	}

	// Construct the target URL
	targetURL := body.URL
	if body.Port != "" {
		targetURL += ":" + body.Port
	}
	targetURL += "/v1/tester"

	log.Infof("Calling target URL: %s", targetURL)

	// Prepare the request payload
	payload := map[string]interface{}{
		"expectedStatus":  body.ExpectedStatus,
		"requestDuration": body.RequestDuration,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to marshal request payload", http.StatusInternalServerError)
		log.WithFields(log.Fields{"err": err}).Error("Failed to marshal request payload")
		return
	}

	// Make the request to the target URL
	resp, err := http.Post(targetURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, "Failed to call target service", http.StatusInternalServerError)
		log.WithFields(log.Fields{"err": err}).Error("Failed to call target service")
		return
	}
	defer resp.Body.Close()

	// Forward the response from the target service
	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to read response body")
	}
	log.Infof("Forwarded response from target service with status code: %d", resp.StatusCode)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	log.Info("Handling /health request")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := make(map[string]string)
	resp["message"] = "Status OK"

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.WithFields(log.Fields{"err": err}).Error("Failed to marshal message")
	}
	log.Info("Health check response sent")
}
