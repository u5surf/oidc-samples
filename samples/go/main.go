package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var host, port, issuer, clientID, clientSecret, redirectURL string
	flag.StringVar(&host, "host", "0.0.0.0", "Listen host")
	flag.StringVar(&port, "port", "8000", "Listen port")
	flag.StringVar(&issuer, "issuer", "https://api.mpin.io", "Backend url")
	flag.StringVar(&clientID, "client-id", "", "OIDC Client Id")
	flag.StringVar(&clientSecret, "client-secret", "", "OIDC Client Secret")
	flag.StringVar(&redirectURL, "redirect-url", "http://localhost:8000/login", "Redirect URL")
	flag.Parse()

	app, err := newApp(issuer, clientID, clientSecret, redirectURL)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", app.index)
	http.HandleFunc("/login", app.login)

	address := host + ":" + port
	log.Printf("Server started. Listening on %v", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
