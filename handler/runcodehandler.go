package handler

import (
	"context"
	"encoding/json"
	"net/http"

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
	jobType, ok := isolate.LanguageNameToJobType[executeCodeRequest.Language]
	if !ok {
		api.Error(rw, http.StatusBadRequest, model.ExecuteCodeError{Error: "invalid job type"})
		return
	}

	languageSpecs := isolate.SupportedLanguageSpecs[executeCodeRequest.ID]
	isolate.Run(context.Background(), jobType, languageSpecs, executeCodeRequest)

	var executeCodeResponse model.ExecuteCodeResponse
	api.Success(rw, http.StatusOK, executeCodeResponse)
}
