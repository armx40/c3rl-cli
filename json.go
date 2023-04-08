package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type malformedRequest struct {
	status int
	msg    string
}

func json_decode_string(data_str string, ret interface{}) error {

	reader := strings.NewReader(data_str)

	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()

	err := dec.Decode(&ret)
	if err != nil {
		// var syntaxError *json.SyntaxError
		// var unmarshalTypeError *json.UnmarshalTypeError
		return err
		// switch {
		// case errors.As(err, &syntaxError):
		// 	msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
		// 	return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// case errors.Is(err, io.ErrUnexpectedEOF):
		// 	msg := fmt.Sprintf("Request body contains badly-formed JSON")
		// 	return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// case errors.As(err, &unmarshalTypeError):
		// 	msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
		// 	return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// case strings.HasPrefix(err.Error(), "json: unknown field "):
		// 	fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		// 	msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
		// 	return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// case errors.Is(err, io.EOF):
		// 	msg := "Request body must not be empty"
		// 	return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		// case err.Error() == "http: request body too large":
		// 	msg := "Request body must not be larger than 1MB"
		// 	return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		// default:
		// 	return err
		// }
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return fmt.Errorf(msg)
	}

	return nil
}
