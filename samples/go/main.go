package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	var addr, issuer, clientID, clientSecret, redirectURL string
	flag.StringVar(&addr, "addr", ":8000", "Listen address")
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

	log.Printf("Server started. Listening on %v", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
