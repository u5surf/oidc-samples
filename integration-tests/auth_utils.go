package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

func pass1(identity Identity, proof []byte, scope ...string) (p1Response *Pass1Response, err error) {
	payload := Pass1Request{
		U:      hex.EncodeToString(proof),
		MPinID: hex.EncodeToString(identity.MPinID),
		DTAs:   identity.DTAs,
		Scope:  scope,
	}

	resp, err := Post(
		"http://api.playground/rps/v2/pass1",
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

func pass2(identity Identity, proof []byte, WID string) (p2Response *Pass2Response, err error) {
	payload := Pass2Request{
		V:      hex.EncodeToString(proof),
		WID:    WID,
		MPinID: hex.EncodeToString(identity.MPinID),
	}

	resp, err := Post(
		"http://api.playground/rps/v2/pass2",
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

func authenticate(authOTT string) (authResponse *AuthenticateResponse, err error) {
	payload := AuthenticateRequest{
		AuthOTT: authOTT,
	}

	resp, err := Post(
		"http://api.playground/rps/v2/authenticate",
		"POST",
		payload,
	)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(resp, &authResponse); err != nil {
		fmt.Println("err", err)
		return nil, err
	}

	return authResponse, nil
}
