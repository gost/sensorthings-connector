package module

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HTTPOperation describes the HTTP operation such as GET POST DELETE.
type HTTPOperation string

// HTTPOperation is a "enumeration" of the HTTP operations needed for all endpoints.
const (
	HTTPOperationGet    HTTPOperation = "GET"
	HTTPOperationPost   HTTPOperation = "POST"
	HTTPOperationPatch  HTTPOperation = "PATCH"
	HTTPOperationDelete HTTPOperation = "DELETE"
)

// EndpointOperation contains the needed information to create an endpoint in the HTTP.Router
type EndpointOperation struct {
	OperationType HTTPOperation     `json:"operation"`
	Path          string            `json:"path"` //relative path to the endpoint for example: /v1.0/myendpoint/
	Handler       httprouter.Handle `json:"-"`
}

// Endpoint holds the information about an endpoint for a module
type Endpoint struct {
	Name       string              `json:"name"`
	Operations []EndpointOperation `json:"operations"`
}

// GetName returns the endpoint name
func (e *Endpoint) GetName() string {
	return e.Name
}

// GetOperations returns the operations for an endpoint
func (e *Endpoint) GetOperations() []EndpointOperation {
	return e.Operations
}

// ErrorResponse is the default response format for sending errors back
type ErrorResponse struct {
	Error ErrorContent `json:"error"`
}

// ErrorContent holds information on the error that occurred
type ErrorContent struct {
	StatusText string `json:"status"`
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

// HandleGetRequest is the default function to handle incoming GET requests
func HandleGetRequest(w http.ResponseWriter, r *http.Request, h *func() interface{}) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	handler := *h
	data := handler()
	SendJSONResponse(w, http.StatusOK, data)
}

// SendJSONResponse sends the desired message to the user
// the message will be marshalled into an indented JSON format
func SendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	b, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		log.Printf("%v", err.Error())
	}

	w.Write(b)
}

// SendError creates an ErrorResponse message and sets it to the user
// using SendJSONResponse
func SendError(w http.ResponseWriter, error error) {
	// Set te status code, default 500 for error, check if there is an ApiError an get
	// the status code
	var statusCode = http.StatusInternalServerError
	if error != nil {
		switch e := error.(type) {
		case EndpointError:
			statusCode = e.GetHTTPErrorStatusCode()
			break
		}
	}

	statusText := http.StatusText(statusCode)
	errorResponse := ErrorResponse{
		Error: ErrorContent{
			StatusText: statusText,
			StatusCode: statusCode,
			Message:    error.Error(),
		},
	}

	SendJSONResponse(w, statusCode, errorResponse)
}
