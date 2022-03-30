package apimodel

// Error models a standard error response.
type Error struct {
	Status      int    `json:"status,omitempty"`      // HTTP status code
	Message     string `json:"message,omitempty"`     // Message describing the problem
	OperationId string `json:"operationId,omitempty"` // ID for identifying relevant logs
}
