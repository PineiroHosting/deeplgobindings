package deeplclient

import (
	"testing"
	"net/http"
	"os"
)

// TestTranslation tests the functionality of the Translate API function within the deeplclient.
func TestTranslation(t *testing.T)  {
	authKey := os.Getenv("DEEPL_TEST_AUTH_KEY")
	if authKey == "" {
		t.Fatal("could not find 'DEEPL_TEST_AUTH_KEY' environment variable")
	}
	client := &Client{
		Client:&http.Client{},
		AuthKey:[]byte(authKey),
	}
	if resp, err :=client.Translate(&TranslationRequest{
		Text:"Hallo Welt!",
		TargetLang:LangEN,
	}); err != nil {
		panic(err)
	} else {
		t.Logf("received response from DeepL API: %+v", resp)
		if len(resp.Translations) != 1 {
			t.Fail()
		}
	}
}
