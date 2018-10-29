package deeplgobindings

import (
	"fmt"
	"strconv"
)

// UnwrappedApiResponseCodeErr represents an error if the API server returns an unexpected status code.
type UnwrappedApiResponseCodeErr int

// Error returns a compact version of all error information in order to implement the error interface.
func (err UnwrappedApiResponseCodeErr) Error() string {
	return fmt.Sprintf("server returned unexpected status code: %d", err)
}

// KnownRequestErrData contains the standard data which can be parsed from json for every error.
type KnownRequestErrData struct {
	// Message holds the error message returned by the server.
	Message string `json:"message"`
}

// WrongRequestErr indicates the response code 400 returned by the remote API server and contains the error message.
type WrongRequestErr struct {
	*KnownRequestErrData
}

// Error returns a compact version of all error information in order to implement the error interface.
func (err *WrongRequestErr) Error() string {
	return fmt.Sprintf("server returned status code 403 (authorization failed): %s", strconv.Quote(err.Message))
}

// AuthFailedErr indicates the response code 403 returned by the remote API server and contains the error message.
// Normally this error occurs if an invalid auth token was provided.
type AuthFailedErr struct {
	*KnownRequestErrData
}

// Error returns a compact version of all error information in order to implement the error interface.
func (err *AuthFailedErr) Error() string {
	return fmt.Sprintf("server returned status code 404 (wrong request): %s", strconv.Quote(err.Message))
}
