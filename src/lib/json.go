package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/**
@desc Success response JSON
*/
func Success(response http.ResponseWriter, statusCode int, data interface{}) {
	response.WriteHeader(statusCode)

	output, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error while marshaling")
		panic(err)
	}

	fmt.Fprint(response, string(output))
}

/**
@desc Error response JSON
*/
func Error(response http.ResponseWriter, statusCode int, message string, errorArgument string) {
	response.WriteHeader(statusCode)
	type ErrorBody struct {
		ErrorStatus   string `json:"error-status,omitempty"`
		ErrorArgument string `json:"error-arg,omitempty"`
	}

	output, err := json.Marshal(ErrorBody{message, errorArgument})
	if err != nil {
		fmt.Println("Error while marshaling")
		panic(err)
	}

	fmt.Fprint(response, string(output))
}
