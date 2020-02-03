package main

import "fmt"

//Generic types

type identity struct {
	MPinID []byte
	Token  []byte
	DTAs   string
}

type header struct {
	Key, Value string
}

type userInfo struct {
	Email string `json:"email"`
}

type cryptoError struct {
	exitCode int
}

func newCryptoError(exitCode int) cryptoError {
	return cryptoError{
		exitCode: exitCode,
	}
}

func (err cryptoError) error() string {
	return fmt.Sprintf("Crypto error exited with code %v", err.exitCode)
}

//Registration responses

type customVerificationResponse struct {
	MPinID        string `json:"mpinId"`
	HashMPinID    string `json:"hashMPinId"`
	ActivateToken string `json:"actToken"`
	ExpireTime    int64  `json:"expireTime"`
}

type registerResponse struct {
	Active     bool   `json:"active"`
	AppID      string `json:"appId"`
	CustomerID string `json:"customerId"`
	ExpireTime int    `json:"expireTime"`
	MPinID     string `json:"mpinId"`
	NowTime    int    `json:"nowTime"`
	RegOTT     string `json:"regOTT"`
}

type signatureResponse struct {
	ClientSecretShare string `json:"clientSecretShare"`
	CS2URL            string `json:"cs2url"`
	Curve             string `json:"curve"`
	DTAs              string `json:"dtas"`
}

type authorizeResponse struct {
	AccessURL string `json:"accessURL"`
	QRURL     string `json:"qrURL"`
	WebOTT    string `json:"webOTT"`
}

type clientSecretResponse struct {
	ClientSecret    string `json:"clientSecret"`
	DVSClientSecret string `json:"dvsClientSecret"`
	CreatedAt       int    `json:"createdAt"`
	Message         string `json:"message"`
	Version         string `json:"version"`
}

//Authentication responses

type pass1Response struct {
	Y string `json:"y"`
}

type pass2Response struct {
	AuthOTT string `json:"authOTT"`
}

type authenticateResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type accessResponse struct {
	Status      string `json:"status"`
	StatusCode  int    `json:"statusCode"`
	UserID      string `json:"userId"`
	RedirectURL string `json:"redirectURL"`
}
