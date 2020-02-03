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
	//all samples generate a new state and redirect us to an /authorize URL, if we're not logged in
	authorizeRequestURL, err := Request(options.sampleURL, "GET", nil)
	if err != nil {
		t.Error("Error making initial request to the home page:", err)
	}

	identity, err := Register(userID, deviceName, pin, string(authorizeRequestURL))
	if err != nil {
		t.Error("Error registering:", err)
	}

	accessResponse, err := Authenticate(identity, userID, pin, string(authorizeRequestURL))
	if err != nil {
		t.Error("Error authenticating:", err)
	}

	//Contains an URL, from which we can fetch the user info
	sessionURL, err := Request(accessResponse.RedirectURL, "GET", nil)
	if err != nil {
		t.Error("Error logging in:", err)
	}

	userInfoResponse, err := Request(string(sessionURL), "GET", nil)
	if err != nil {
		t.Error("Error getting user info:", err)
	}

	var userInfo userInfo
	if err = json.Unmarshal(userInfoResponse, &userInfo); err != nil {
		t.Error("Failed to unmarshal user info")
	}

	if userInfo.Email != userID {
		t.Error("UserID  mismatch")
	}
}