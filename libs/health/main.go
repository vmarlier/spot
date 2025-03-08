package health

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := make(map[string]string)
	resp["message"] = "Status OK"

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal message")
	}
	return
}
