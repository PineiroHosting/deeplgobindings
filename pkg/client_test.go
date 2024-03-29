package deeplclient

import (
	"net/http"
	"os"
	"testing"
)

const endpointUrl = "https://api-free.deepl.com/v2/"

func TestGetUsage(t *testing.T) {
	authKey := os.Getenv("DEEPL_TEST_AUTH_KEY")
	if authKey == "" {
		t.Fatal("could not find 'DEEPL_TEST_AUTH_KEY' environment variable")
	}
	client := &Client{
		Client:      &http.Client{},
		AuthKey:     []byte(authKey),
		EndpointUrl: endpointUrl,
	}
	if resp, err := client.GetUsage(); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("received usage response from DeepL API: %+v", resp)
		if resp.CharacterCount < 0 {
			t.Fail()
		}
		if resp.CharacterLimit < 0 {
			t.Fail()
		}
	}
}
