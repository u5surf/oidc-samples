package main

import (
	"flag"
	"os"
	"testing"
)

var options struct {
	clientID     string
	clientSecret string
	sampleHost   string
	mpinHost     string
}

func TestMain(m *testing.M) {
	flag.StringVar(&options.clientSecret, "client-secret", "", "the client secret for our portal app")
	flag.StringVar(&options.clientID, "client-id", "", "the client id for our portal app")
	flag.StringVar(&options.sampleHost, "sample-host", "", "the host for the sample")
	flag.StringVar(&options.mpinHost, "mpin-host", "", "the host for mpin")

	flag.Parse()

	os.Exit(m.Run())
}
