package model

// ErrorDto models a standard error response.
type ErrorDto struct {
	Status      int    `json:"status,omitempty"`      // HTTP status code
	Message     string `json:"message,omitempty"`     // Message describing the problem
	OperationId string `json:"operationId,omitempty"` // ID for identifying relevant logs
}
