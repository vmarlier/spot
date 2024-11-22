package main

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {

	addr := ":8080"
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)

	log.Infof("dummy-test is listening at %s...", addr)
	log.Error(http.ListenAndServe(addr, mux))
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
