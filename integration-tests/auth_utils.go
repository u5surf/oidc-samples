package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"

	"code.miracl.com/mfa/pkg/gomiracl"
	"code.miracl.com/mfa/pkg/gomiracl/bindings"
	"code.miracl.com/mfa/pkg/gomiracl/wrap"
)

func Authenticate(identity identity, userID string, pin int, authorizeRequestURL string) (accessResponse *accessResponse, err error) {
	//we'll need data from this response (wid) for pass2 request
	//otherwise, for some reason, if wid is not passed, the authOTT, that you'll need later is not stored in redis
	authorizeResponse, err := authorizeRequest(authorizeRequestURL)
	if err != nil {
		return nil, err
	}

	//Get pass1 proof from the token and pin (this is the U param in /pass1)
	rand := bindings.NewRand([]byte{})
	X := make([]byte, 32)
	proof := make([]byte, 65)
	xR, S, U, _, err := wrap.Client1BN254CX(int(gomiracl.SHA256), 0, identity.MPinID, rand, X, pin, identity.Token, proof)
	if err != nil {
		return nil, err
	}

	//rps/v2/pass1
	p1Response, err := pass1Request(identity, U, "oidc")

	//Get V (used in /pass2) param using Y param from the pass1 response
	V, err := wrap.Client2BN254CX(xR, hex2bytes(p1Response.Y), S)
	if err != nil {
		return nil, err
	}

	//rps/v2/pass2
	qrURL, err := url.Parse(authorizeResponse.QRURL)
	if err != nil {
		return nil, err
	}
	p2Response, err := pass2Request(identity, V, qrURL.Fragment)
	if err != nil {
		return nil, err
	}

	//rps/v2/authenticate
	authResponse, err := authenticateRequest(p2Response.AuthOTT)
	if err != nil {
		return nil, err
	}
	if authResponse.Status != 200 {
		return nil, fmt.Errorf(authResponse.Message)
	}

	//rps/v2/access
	accessResponse, err = accessRequest(authorizeResponse.WebOTT)
	if err != nil {
		return nil, err
	}

	return accessResponse, nil
}

func pass1Request(identity identity, proof []byte, scope ...string) (p1Response *pass1Response, err error) {
	payload := struct {
		U      string   `json:"U"`
		MPinID string   `json:"mpin_id"`
		DTAs   string   `json:"dtas"`
		Scope  []string `json:"scope"`
	}{
		hex.EncodeToString(proof),
		hex.EncodeToString(identity.MPinID),
		identity.DTAs,
		scope,
	}

	resp, err := Request(
		options.apiURL+"/rps/v2/pass1",
		"POST",
		payload,
	)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resp, &p1Response); err != nil {
		return nil, err
	}

	return p1Response, nil
}

func pass2Request(identity identity, proof []byte, WID string) (p2Response *pass2Response, err error) {
	payload := struct {
		V      string `json:"V"`
		WID    string `json:"WID"`
		MPinID string `json:"mpin_id"`
	}{
		hex.EncodeToString(proof),
		WID,
		hex.EncodeToString(identity.MPinID),
	}

	resp, err := Request(
		options.apiURL+"/rps/v2/pass2",
		"POST",
		payload,
	)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resp, &p2Response); err != nil {
		return nil, err
	}

	return p2Response, nil
}

func authenticateRequest(authOTT string) (authResponse *authenticateResponse, err error) {
	payload := struct {
		AuthOTT string `json:"authOTT"`
	}{
		authOTT,
	}

	resp, err := Request(
		options.apiURL+"/rps/v2/authenticate",
		"POST",
		payload,
	)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resp, &authResponse); err != nil {
		return nil, err
	}

	return authResponse, nil
}

func accessRequest(webOTT string) (accessResponse *accessResponse, err error) {
	payload := struct {
		WebOTT string `json:"webOTT"`
	}{
		webOTT,
	}

	resp, err := Request(
		options.apiURL+"/rps/v2/access",
		"POST",
		payload,
	)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resp, &accessResponse); err != nil {
		return nil, err
	}

	return accessResponse, nil
}
