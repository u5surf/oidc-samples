package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"

	"code.miracl.com/mfa/pkg/gomiracl"
	"code.miracl.com/mfa/pkg/gomiracl/wrap"
)

func register(userID string, deviceName string, pin int, authorizeRequestURL string) (i identity, err error) {

	// Call to /authorize endpoint.
	authorizeResponse, err := authorizeRequest(authorizeRequestURL)
	if err != nil {
		return identity{}, err
	}

	// Call to /activate/initiate endpoint.
	cvResponse, err := customVerifyRequest(userID, deviceName)
	if err != nil {
		return identity{}, err
	}

	// Call to /rps/v2/user endpoint.
	qrURL, err := url.Parse(authorizeResponse.QRURL)
	if err != nil {
		return identity{}, err
	}
	regResponse, err := registerRequest(userID, deviceName, qrURL.Fragment, cvResponse.ActivateToken)
	if err != nil {
		return identity{}, err
	}

	// Call to /signature endpoint.
	sigResponse, err := signatureRequest(regResponse.MPinID, regResponse.RegOTT)
	if err != nil {
		return identity{}, err
	}

	// Call to /dta/ID endpoint.
	csResponse, err := clientSecretRequest(sigResponse.CS2URL)

	// Combine both client secrets.
	Q, err := wrap.RecombineG1BN254CX(hex2bytes(sigResponse.ClientSecretShare), hex2bytes(csResponse.ClientSecret))
	if err != nil {
		return identity{}, err
	}

	// First extract pin from the combine client secret, in order to get the token.
	CS, err := wrap.ExtractPINBN254CX(int(gomiracl.SHA256), hex2bytes(regResponse.MPinID), pin, Q)
	if err != nil {
		return identity{}, err
	}

	return identity{
		MPinID: hex2bytes(regResponse.MPinID),
		Token:  CS,
		DTAs:   sigResponse.DTAs,
	}, nil
}

func authorizeRequest(requestURL string) (*authorizeResponse, error) {
	resp, err := request(
		requestURL,
		"POST",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var authorizeResponse *authorizeResponse
	if err := json.Unmarshal(resp, &authorizeResponse); err != nil {
		return nil, err
	}

	return authorizeResponse, nil
}

func customVerifyRequest(userID string, deviceName string) (*customVerificationResponse, error) {
	payload := struct {
		UserID     string `json:"userId"`
		DeviceName string `json:"deviceName"`
	}{
		userID,
		deviceName,
	}

	clientIdAndSecret := options.clientID + ":" + options.clientSecret
	authHeaderValue := "Basic " + b64.StdEncoding.EncodeToString([]byte(clientIdAndSecret))

	resp, err := request(
		options.apiURL+"/activate/initiate",
		"POST",
		payload,
		header{Key: "Authorization", Value: authHeaderValue},
	)
	if err != nil {
		return nil, err
	}

	var customVerificationResponse *customVerificationResponse
	if err := json.Unmarshal(resp, &customVerificationResponse); err != nil {
		return nil, err
	}

	return customVerificationResponse, nil
}

func registerRequest(userID string, deviceName string, wid string, activateCode string) (*registerResponse, error) {
	payload := struct {
		DeviceName   string `json:"deviceName"`
		UserID       string `json:"userId"`
		WID          string `json:"wid"`
		ActivateCode string `json:"activateCode"`
	}{
		DeviceName:   deviceName,
		UserID:       userID,
		WID:          wid,
		ActivateCode: activateCode,
	}
	resp, err := request(
		options.apiURL+"/rps/v2/user",
		"PUT",
		payload,
		header{Key: "X-MIRACL-CID", Value: "mcl"},
	)
	if err != nil {
		return nil, err
	}

	var registerResponse *registerResponse
	if err := json.Unmarshal(resp, &registerResponse); err != nil {
		return nil, err
	}

	return registerResponse, nil
}

func signatureRequest(mpinID string, regOTT string) (*signatureResponse, error) {
	resp, err := request(
		fmt.Sprintf(options.apiURL+"/rps/v2/signature/%v?regOTT=%v", mpinID, regOTT),
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var sigResponse *signatureResponse
	if err := json.Unmarshal(resp, &sigResponse); err != nil {
		return nil, err
	}

	return sigResponse, nil
}

func clientSecretRequest(cs2url string) (*clientSecretResponse, error) {
	resp, err := request(
		cs2url,
		"GET",
		nil,
	)
	if err != nil {
		return nil, err
	}

	var csResponse *clientSecretResponse
	if err := json.Unmarshal(resp, &csResponse); err != nil {
		return nil, err
	}

	return csResponse, nil
}
