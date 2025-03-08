package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/vmarlier/spot/internal/api/health"
)

func main() {
	addr := ":8080"
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health.Handler)

	logrus.Infof("server is listening at %s...", addr)
	logrus.Error(http.ListenAndServe(addr, mux))
}
