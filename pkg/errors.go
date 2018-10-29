package deeplgobindings

import (
	"fmt"
)

// UnwrappedApiResponseCodeErr represents an error if the API server returns an unexpected status code.
type UnwrappedApiResponseCodeErr int

// Error returns a compact version of all error information in order to implement the error interface.
func (err UnwrappedApiResponseCodeErr) Error() string {
	return fmt.Sprintf("server returned unexpected status code: %d", err)
}
