package deeplclient

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
	return fmt.Sprintf("server returned status code 400 (wrong request): %s", strconv.Quote(err.Message))
}

// AuthFailedErr indicates the response code 403 returned by the remote API server and contains the error message.
// Normally this error occurs if an invalid auth token was provided.
type AuthFailedErr struct {
	*KnownRequestErrData
}

// Error returns a compact version of all error information in order to implement the error interface.
func (err *AuthFailedErr) Error() string {
	return fmt.Sprintf("server returned status code 403 (autorization failed): %s", strconv.Quote(err.Message))
}

// RequestEntityTooLargeErr indicates the response code 413 returned by the remote API server and contains the error message.
// Normally this error occurs if the request size exceeds the current limit.
type RequestEntityTooLargeErr struct {
	*KnownRequestErrData
}

// Error returns a compact version of all error information in order to implement the error interface.
func (err *RequestEntityTooLargeErr) Error() string {
	return fmt.Sprintf("server returned status code 413 (request entity too large): %s", strconv.Quote(err.Message))
}

// TooManyRequestsErr indicates the response code 429 returned by the remote API server and contains the error message.
// Normally this error occurs if too many requests have been sent in a short amount of time.
type TooManyRequestsErr struct {
	*KnownRequestErrData
}

// Error returns a compact version of all error information in order to implement the error interface.
func (err *TooManyRequestsErr) Error() string {
	return fmt.Sprintf("server returned status code 429 (too many request): %s", strconv.Quote(err.Message))
}

// AuthFailedErr indicates the response code 403 returned by the remote API server and contains the error message.
// Normally this error occurs if the character limit has been reached.
type QuotaExceededErr struct {
	*KnownRequestErrData
}

// Error returns a compact version of all error information in order to implement the error interface.
func (err *QuotaExceededErr) Error() string {
	return fmt.Sprintf("server returned status code 456 (quota exceeded): %s", strconv.Quote(err.Message))
}
