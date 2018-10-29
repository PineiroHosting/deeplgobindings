package main

import (
	"flag"
	"fmt"
	"github.com/PineiroHosting/deeplgobindings/pkg"
	"net/http"
)

func main() {
	var authKey = flag.String("authkey", "", "DeepL developer plan API auth key in order to access the API.")
	flag.Parse()
	if *authKey == "" {
		fmt.Println("Please provide a valid auth key!")
		flag.PrintDefaults()
		return
	}
	client := &deeplclient.Client{
		AuthKey: []byte(*authKey),
		Client:  &http.Client{},
	}
	resp, err := client.GetUsage()
	// basic error handling because it is an example
	if err != nil {
		panic(err)
	}
	fmt.Println("===================================================")
	fmt.Printf("| Monthly usage: %32s |\n", fmt.Sprintf("%f%s (%d/%d)", float64(resp.CharacterCount) / float64(resp.CharacterLimit) / 100.0, "%", resp.CharacterCount, resp.CharacterLimit))
	fmt.Println("===================================================")
}
