package deeplgobindings

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

// doApiFunction is an internally used function to execute API functions more easily. The param uri should not begin with
// a slash character.
func (client *DeeplClient) doApiFunction(uri string, values *url.Values) (resp *http.Response, err error) {
	// add authentication header and encode values
	values.Set(authKeyparamName, string(client.AuthKey))
	body := strings.NewReader(values.Encode())
	// create new http request
	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, fmt.Sprintf("%s%s", deeplBaseApiUrl, uri), body); err != nil {
		return
	}
	req.ContentLength = body.Size()
	// add header to allow the server to identify the POST request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if resp, err = client.Do(req); err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		return
	default:
		err = UnwrappedApiResponseCodeErr(resp.StatusCode)
		return nil, err
	}
	return
}
