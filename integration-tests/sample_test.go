package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
)

func TestAuth(t *testing.T) {
	name := uuid.New().String()
	userID := fmt.Sprintf("%v@example.com", name)
	deviceName := "The device of " + name
	pin := randPIN()

	// All samples generate a new state and redirect us to an /authorize URL, if we're not logged in.
	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	client := newSampleClient(options.sampleURL, httpClient)

	authorizeRequestURL, err := client.authorize()
	if err != nil {
		t.Fatal("Error making initial request to the home page:", err)
	}

	identity, err := register(httpClient, userID, deviceName, pin, string(authorizeRequestURL))
	if err != nil {
		t.Fatal("Error registering:", err)
	}

	accessResponse, err := authenticate(httpClient, identity, userID, pin, string(authorizeRequestURL))
	if err != nil {
		t.Fatal("Error authenticating:", err)
	}

	// Contains an URL, from which we can fetch the user info.
	err = client.login(accessResponse.RedirectURL)
	if err != nil {
		t.Fatal("Error logging in:", err)
	}

	userInfo, err := client.getUserInfo()
	if err != nil {
		t.Fatal("Error getting user info:", err)
	}

	if userInfo.Email != userID {
		t.Fatal("UserID  mismatch")
	}
}
