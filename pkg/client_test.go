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
		t.Fatal(err)
	} else {
		t.Logf("received translation response from DeepL API: %+v", resp)
		if len(resp.Translations) != 1 {
			t.Fail()
		}
	}
}

func TestGetUsage(t *testing.T) {
	authKey := os.Getenv("DEEPL_TEST_AUTH_KEY")
	if authKey == "" {
		t.Fatal("could not find 'DEEPL_TEST_AUTH_KEY' environment variable")
	}
	client := &Client{
		Client:&http.Client{},
		AuthKey:[]byte(authKey),
	}
	if resp, err :=client.GetUsage(); err != nil {
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