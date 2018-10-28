package deeplgobindings

import (
	"fmt"
	"io"
	"net/http"
)

const (
	deeplBaseApiUrl  = "https://api.deepl.com/v2/"
	authKeyparamName = "auth_key"
)

// DeeplClient allows easy access to the DeepL API by providing methods for each API function.
type DeeplClient struct {
	*http.Client
	// AuthKey stores the authentication key required to get access DeepL's API.
	AuthKey []byte
}

func (client *DeeplClient) doApiFunction(uri string, body io.Reader) (resp *http.Response, err error) {
	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", deeplBaseApiUrl, uri), body); err != nil {
		return
	}
	if resp, err = client.Do(req); err != nil {
		return nil, err
	}
	return
}
