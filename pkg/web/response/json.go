package response

import (
	"encoding/json"
	"net/http"

	"github.com/scraly/go.common/pkg/log"
)

// JSON writes JSON to an http.ResponseWriter with the corresponding status code
func JSON(w http.ResponseWriter, status int, data interface{}) {
	// Get rid of the invalid status codes
	if status < 100 || status > 599 {
		status = 200
	}

	// Try to marshal the input
	result, err := json.Marshal(data)
	if err != nil {
		// Set the result to the default value to prevent empty responses
		result = []byte(`{"status":500,"message":"Error occurred while marshalling the response body"}`)
	}

	// Set the response's content type to JSON
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// Write the result
	w.WriteHeader(status)

	// Write body
	_, err = w.Write(result)
	log.CheckErr("Unable to write JSON data", err)
}
