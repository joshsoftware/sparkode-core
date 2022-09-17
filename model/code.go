package model

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
