package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestAuth(t *testing.T) {
	name := uuid.New().String()
	userID := fmt.Sprintf("%v@example.com", name)
	deviceName := "The device of " + name
	pin := randPIN()

	// All samples generate a new state and redirect us to an /authorize URL, if we're not logged in.
	authorizeRequestURL, err := request(options.sampleURL, "GET", nil)
	if err != nil {
		t.Fatal("Error making initial request to the home page:", err)
	}

	identity, err := register(userID, deviceName, pin, string(authorizeRequestURL))
	if err != nil {
		t.Fatal("Error registering:", err)
	}

	accessResponse, err := authenticate(identity, userID, pin, string(authorizeRequestURL))
	if err != nil {
		t.Fatal("Error authenticating:", err)
	}

	// Contains an URL, from which we can fetch the user info.
	sessionURL, err := request(accessResponse.RedirectURL, "GET", nil)
	if err != nil {
		t.Fatal("Error logging in:", err)
	}

	userInfoResponse, err := request(string(sessionURL), "GET", nil)
	if err != nil {
		t.Fatal("Error getting user info:", err)
	}

	var userInfo userInfo
	if err = json.Unmarshal(userInfoResponse, &userInfo); err != nil {
		t.Fatal("Failed to unmarshal user info")
	}

	if userInfo.Email != userID {
		t.Fatal("UserID  mismatch")
	}
}
