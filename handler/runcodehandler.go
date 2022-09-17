package handler

import (
	"encoding/json"
	"net/http"

	"github.com/joshsoftware/sparkode-core/api"
	"github.com/joshsoftware/sparkode-core/service"
)

type PingResponse struct {
	Message string `json:"message"`
}

type ExecuteCodeRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	Input    string `json:"input"`
}

type ExecuteCodeResponse struct {
	Status    bool    `json:"status"`
	Output    string  `json:"output"`
	TimeTaken float32 `json:"time_taken"`
}

type ExecuteCodeError struct {
	Error string `json:"error"`
}

func PingHandler(rw http.ResponseWriter, req *http.Request) {
	pingResponse := PingResponse{Message: "pong"}
	api.Success(rw, http.StatusOK, pingResponse)
}

func RuncodeHandler(service service.Service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var executeCodeRequest ExecuteCodeRequest
		err := json.NewDecoder(req.Body).Decode(&executeCodeRequest)
		if err != nil {
			api.Error(rw, http.StatusBadRequest, ExecuteCodeError{Error: "incorrect request body"})
			return
		}

		if executeCodeRequest.Language == "" || executeCodeRequest.Code == "" {
			api.Error(rw, http.StatusBadRequest, ExecuteCodeError{Error: "language/program field must not be empty"})
			return
		}
		res, err := service.Run(req.Context(), executeCodeRequest.Code, executeCodeRequest.Language, executeCodeRequest.Input)
		var executeCodeResponse ExecuteCodeResponse
		executeCodeResponse.Output = res
		api.Success(rw, http.StatusOK, executeCodeResponse)
	})
}
