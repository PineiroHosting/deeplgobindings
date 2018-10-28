package deeplgobindings

import (
	"net/http"
)

const (
	AuthKeyparamName = "auth_key"
)

// DeeplClient allows easy access to the DeepL API by providing methods for each API function.
type DeeplClient struct {
	*http.Client
	// AuthKey stores the authentication key required to get access DeepL's API.
	AuthKey []byte
}
