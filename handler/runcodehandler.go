package handler

import (
	"encoding/json"
	"net/http"

	logger "github.com/sirupsen/logrus"
)

type RuncodeResponse struct {
	Output string `json:"output"`
}

func RuncodeHandler(rw http.ResponseWriter, req *http.Request) {
	response := RuncodeResponse{Output: "output"}

	respBytes, err := json.Marshal(response)
	if err != nil {
		logger.WithField("err", err.Error()).Error("Error marshalling ping response")
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.Write(respBytes)
}
