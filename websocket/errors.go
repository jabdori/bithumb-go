package websocket

// WebSocket error type constants as defined in Bithumb API documentation
const (
	WSErrWrongFormat        = "WRONG_FORMAT"
	WSErrNoTicket           = "NO_TICKET"
	WSErrNoType             = "NO_TYPE"
	WSErrNoCodes            = "NO_CODES"
	WSErrInvalidParam       = "INVALID_PARAM"
	WSErrInvalidParamFormat = "INVALID_PARAM_FORMAT"
)

// WSError represents a WebSocket API error response
type WSError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

// WSErrorResponse represents the full error response structure
type WSErrorResponse struct {
	Error WSError `json:"error"`
}
