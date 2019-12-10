package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func authorize(userID string, clientID string, redirectURI string, nonce string, scope ...string) (*AuthorizeResponse, error) {
	query := url.Values{
		"client_id":     []string{clientID},
		"redirect_uri":  []string{redirectURI},
		"response_type": []string{"code"},
		"scope":         []string{strings.Join(scope, " ")},
		"nonce":         []string{nonce},
	}.Encode()

	resp, err := Post(
		"https://api.mpin.io/authorize"+"?"+query,
		"POST",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var authorizeResponse *AuthorizeResponse
	if err := json.Unmarshal(resp, &authorizeResponse); err != nil {
		return nil, err
	}

	return authorizeResponse, nil
}

func customVerify(userID string, deviceName string) (*CustomVerificationResponse, error) {
	payload := CustomVerificationRequest{
		UserID:     userID,
		DeviceName: deviceName,
	}

	clientIdAndSecret := "l2tfcf23rxl3h:6lP2dBrlQHVngFGQzaDW7ghdrLUJOxaDC74E68H15gI"
	authHeaderValue := "Basic " + b64.StdEncoding.EncodeToString([]byte(clientIdAndSecret))

	resp, err := Post(
		"https://api.mpin.io/activate/initiate",
		"POST",
		payload,
		Header{Key: "Authorization", Value: authHeaderValue},
	)
	if err != nil {
		return nil, err
	}

	var customVerificationResponse *CustomVerificationResponse
	if err := json.Unmarshal(resp, &customVerificationResponse); err != nil {
		return nil, err
	}

	return customVerificationResponse, nil
}

func register(userID string, deviceName string, wid string, activateCode string) (*RegisterResponse, error) {
	payload := RegisterRequest{
		DeviceName:   deviceName,
		UserID:       userID,
		WID:          wid,
		ActivateCode: activateCode,
	}
	resp, err := Post(
		"https://api.mpin.io/rps/v2/user",
		"PUT",
		payload,
		Header{Key: "X-MIRACL-CID", Value: "mcl"},
	)
	if err != nil {
		return nil, err
	}

	var registerResponse *RegisterResponse
	if err := json.Unmarshal(resp, &registerResponse); err != nil {
		return nil, err
	}

	return registerResponse, nil
}

func signature(mpinID string, regOTT string) (*SignatureResponse, error) {
	resp, err := Post(
		fmt.Sprintf("https://api.mpin.io/rps/v2/signature/%v?regOTT=%v", mpinID, regOTT),
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var sigResponse *SignatureResponse
	if err := json.Unmarshal(resp, &sigResponse); err != nil {
		return nil, err
	}

	return sigResponse, nil
}

func clientSecret(cs2url string) (*ClientSecretResponse, error) {
	resp, err := Post(
		cs2url,
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var csResponse *ClientSecretResponse
	if err := json.Unmarshal(resp, &csResponse); err != nil {
		return nil, err
	}

	return csResponse, nil
}
