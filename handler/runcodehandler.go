package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/joshsoftware/sparkode-core/api"
	"github.com/joshsoftware/sparkode-core/isolate"
	"github.com/joshsoftware/sparkode-core/model"
)

func PingHandler(rw http.ResponseWriter, req *http.Request) {
	pingResponse := model.PingResponse{Message: "pong"}
	api.Success(rw, http.StatusOK, pingResponse)
}

func RuncodeHandler(rw http.ResponseWriter, req *http.Request) {
	var executeCodeRequest model.ExecuteCodeRequest
	err := json.NewDecoder(req.Body).Decode(&executeCodeRequest)
	if err != nil {
		api.Error(rw, http.StatusBadRequest, model.ExecuteCodeError{Error: "incorrect request body"})
		return
	}

	if executeCodeRequest.Language == "" || executeCodeRequest.Code == "" {
		api.Error(rw, http.StatusBadRequest, model.ExecuteCodeError{Error: "language/program field must not be empty"})
		return
	}

	start := time.Now()
	languageSpecs, ok := isolate.SupportedLanguageSpecs[executeCodeRequest.ID]
	if !ok {
		api.Error(rw, http.StatusBadRequest, model.ExecuteCodeError{Error: "requested language doesnot exist"})
		return
	}

	stdout, stderr, err := isolate.Run(context.Background(), languageSpecs, executeCodeRequest)
	if err != nil {
		api.Error(rw, http.StatusInternalServerError, model.ExecuteCodeError{Error: err.Error()})
		return
	}

	timeTaken := time.Since(start)
	var executeCodeResponse model.ExecuteCodeResponse

	executeCodeResponse.TimeTaken = float32(timeTaken.Seconds())

	if stderr != "" {
		executeCodeResponse.Output = ""
		executeCodeResponse.Status = false
		api.Success(rw, http.StatusOK, executeCodeResponse)
	}

	executeCodeResponse.Status = true
	executeCodeResponse.Output = stdout

	api.Success(rw, http.StatusOK, executeCodeResponse)
}
