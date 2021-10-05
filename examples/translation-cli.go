//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/PineiroHosting/deeplgobindings/pkg"
	"net/http"
	"net/url"
	"os"
)

func main() {
	var authKey = flag.String("authkey", "", "DeepL developer plan API auth key in order to access the API.")
	var endpointUrl = flag.String("endpoint", "https://api.deepl.com/v2/", "DeepL endpoint url (use https://api-free.deepl.com/v2/ for the free-api)")
	flag.Parse()
	if *authKey == "" {
		fmt.Println("Please provide a valid auth key!")
		flag.PrintDefaults()
		return
	}
	if _, err := url.Parse(*endpointUrl); err != nil {
		fmt.Println("Please provide a valid endpoint url!")
		flag.PrintDefaults()
		return
	}
	client := &deeplclient.Client{
		AuthKey:     []byte(*authKey),
		Client:      &http.Client{},
		EndpointUrl: *endpointUrl,
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter text which should be translated to English. Enter 'stop' to stop.")
	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "stop" {
			break
		}
		resp, err := client.Translate(&deeplclient.TranslationRequest{
			Text:       text,
			TargetLang: deeplclient.LangEN,
		})
		// basic error handling because it is an example
		if err != nil {
			panic(err)
		}
		fmt.Printf("[%s->EN] %s\n", resp.Translations[0].DetectedSourceLanguage, resp.Translations[0].Text)
		fmt.Print("> ")
	}
	fmt.Println("Bye!")
}
