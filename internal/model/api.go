package model

type HealthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
