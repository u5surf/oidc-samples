package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	oidc "github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

type app struct {
	provider    *oidc.Provider
	verifier    *oidc.IDTokenVerifier
	oauthConfig *oauth2.Config

	stateStorage   stateStorage    // This is sample implementation. Don't do it in production
	sessionStorage userInfoStorage // This is sample implementation. Don't do it in production

	random func() string
}

func newApp(issuer, clientID, clientSecret, redirectURL, proxyHost, proxyPort string) (*app, error) {
	ctx := context.Background()

	if proxyHost != "" && proxyPort != "" {
		proxyUrl, err := url.Parse(fmt.Sprintf("http://%v:%v", proxyHost, proxyPort))
		if err != nil {
			return nil, err
		}

		http.DefaultTransport = &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

		myClient := &http.Client{}
		ctx = oidc.ClientContext(context.Background(), myClient)
	}

	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, err
	}

	return &app{
		provider: provider,
		verifier: provider.Verifier(&oidc.Config{ClientID: clientID}),
		oauthConfig: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  redirectURL,
			Scopes:       []string{oidc.ScopeOpenID, "email"},
		},
		stateStorage:   stateStorage{},
		sessionStorage: userInfoStorage{},
	}, nil
}

func (a *app) index(w http.ResponseWriter, r *http.Request) {
	loginRedirect := func() {
		http.Redirect(w, r, a.oauthConfig.AuthCodeURL(a.stateStorage.new()), http.StatusFound)
	}

	sessionID := r.URL.Query().Get("session")
	if sessionID == "" {
		loginRedirect()
		return
	}

	userInfo := a.sessionStorage.get(sessionID)
	if userInfo == nil {
		loginRedirect()
		return
	}

	response, err := json.Marshal(&userInfo)
	if err != nil {
		log.Println("Error marshaling UserInfo: %w", err)
		http.Error(w, "Error marshaling UserInfo.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (a *app) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// validate the state
	state := r.URL.Query().Get("state")
	if !a.stateStorage.isValid(state) {
		http.Error(w, "Invalid state.", http.StatusBadRequest)
		return
	}
	a.stateStorage.pop(state)

	// exchange the code
	oauth2Token, err := a.oauthConfig.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		log.Println("Error exchanging token: %w", err)
		http.Error(w, "Exchanging code failed.", http.StatusInternalServerError)
		return
	}

	// verify the token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		log.Printf("Error getting id token")
		http.Error(w, "Invalid token.", http.StatusInternalServerError)
		return
	}
	_, err = a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Printf("ID Token verification failed: %v", err)
		http.Error(w, "Invalid ID token.", http.StatusBadRequest)
		return
	}

	// get the user info from the provider
	userInfo, err := a.provider.UserInfo(ctx, oauth2.StaticTokenSource(oauth2Token))
	if err != nil {
		log.Printf("Error getting UserInfo: %v", err)
		http.Error(w, "UserInfo request failed.", http.StatusInternalServerError)
		return
	}

	sessionID := a.sessionStorage.add(userInfo)
	http.Redirect(w, r, fmt.Sprintf("/?session=%v", sessionID), http.StatusFound)
}
