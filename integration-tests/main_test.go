package main

import (
	"flag"
	"fmt"
	"os"
	"testing"
)

var options struct {
	clientID     string
	clientSecret string
	redirectURL  string
	apiURL       string
	sampleURL    string
}

func TestMain(m *testing.M) {
	flag.StringVar(&options.clientSecret, "client-secret", "", "the client secret for the portal app")
	flag.StringVar(&options.clientID, "client-id", "", "the client id for the portal app")
	flag.StringVar(&options.redirectURL, "redirect-url", "http://localhost:8000/login", "the redirect url from the portal app")
	flag.StringVar(&options.apiURL, "api-url", "https://api.mpin.io", "the mpin api URL")
	flag.StringVar(&options.sampleURL, "sample-url", "http://localhost:8000", "the sample URL")

	flag.Parse()

	if options.clientSecret == "" && options.clientID == "" {
		fmt.Println("ERROR: client-id and client-secret args are missing.\nUse -h flag to see all args.")
		os.Exit(1)
	}

	os.Exit(m.Run())
}
